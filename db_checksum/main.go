// Package main demonstrates how database page checksums work in PostgreSQL and MySQL/InnoDB.
// This is an educational demo showing corruption detection mechanisms.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// PostgreSQL uses 8KB pages, MySQL/InnoDB uses 16KB pages
const (
	PostgreSQLPageSize = 8192  // 8KB
	InnoDBPageSize     = 16384 // 16KB
)

// PostgreSQL Page Header (simplified)
// Reference: src/include/storage/bufpage.h
type PostgreSQLPageHeader struct {
	PDLsnXLogRecPtr uint64 // LSN: next byte after last byte of xlog record
	PDChecksum      uint16 // Page checksum
	PDFlags         uint16 // Flag bits
	PDLower         uint16 // Offset to start of free space
	PDUpper         uint16 // Offset to end of free space
	PDSpecial       uint16 // Offset to start of special space
	PDPageSizeVer   uint16 // Page size and layout version number
	PDPruneXid      uint32 // Oldest prunable XID, or zero if none
}

// InnoDB Page Header (simplified)
// Reference: storage/innobase/include/fil0fil.h
type InnoDBPageHeader struct {
	Checksum   uint32 // 0-3: Checksum
	PageNumber uint32 // 4-7: Page number
	PrevPage   uint32 // 8-11: Previous page
	NextPage   uint32 // 12-15: Next page
	LSN        uint64 // 16-23: Log sequence number
	PageType   uint16 // 24-25: Page type
	FlushLSN   uint64 // 26-33: Flush LSN (only in first page)
	SpaceID    uint32 // 34-37: Space ID
}

// InnoDB Page Trailer
type InnoDBPageTrailer struct {
	OldChecksum uint32 // Old-style checksum
	LSNLow      uint32 // Low 4 bytes of LSN
}

// PostgreSQLPage simulates a PostgreSQL data page
type PostgreSQLPage struct {
	Header PostgreSQLPageHeader
	Data   []byte // User data area
}

// InnoDBPage simulates an InnoDB data page
type InnoDBPage struct {
	Header  InnoDBPageHeader
	Data    []byte // User data area (16384 - 38 - 8 = 16338 bytes)
	Trailer InnoDBPageTrailer
}

// ============================================================================
// PostgreSQL Checksum Implementation
// Based on: src/port/pg_crc32c.c and src/backend/storage/page/checksum.c
// PostgreSQL uses a FNV-1a variant optimized for SIMD operations
// ============================================================================

// pgChecksumBlock calculates PostgreSQL-style checksum
// This is a simplified version of pg_checksum_page()
func pgChecksumBlock(data []byte, blockNumber uint32) uint16 {
	// PostgreSQL uses FNV-1a hash with SIMD optimization
	// We simulate a simplified version here

	const fnvPrime = 0x01000193
	const fnvOffsetBasis = 0x811c9dc5

	hash := uint32(fnvOffsetBasis)

	// Mix in the block number
	hash ^= blockNumber
	hash *= fnvPrime

	// Process data in 4-byte chunks (simplified from actual 128-bit SIMD)
	for i := 0; i < len(data)-3; i += 4 {
		// Skip the checksum field itself (bytes 8-9 in header)
		if i == 8 {
			continue
		}
		val := binary.LittleEndian.Uint32(data[i:])
		hash ^= val
		hash *= fnvPrime
	}

	// Fold 32-bit to 16-bit
	return uint16(hash) ^ uint16(hash>>16)
}

// NewPostgreSQLPage creates a new PostgreSQL page with sample data
func NewPostgreSQLPage(blockNumber uint32, userData string) *PostgreSQLPage {
	page := &PostgreSQLPage{
		Header: PostgreSQLPageHeader{
			PDLsnXLogRecPtr: 0x0000000100000001,
			PDFlags:         0,
			PDLower:         24, // After header
			PDUpper:         8192,
			PDSpecial:       8192,
			PDPageSizeVer:   0x2004, // 8192 bytes, version 4
			PDPruneXid:      0,
		},
		Data: make([]byte, PostgreSQLPageSize-24), // Subtract header size
	}

	// Write user data
	copy(page.Data, userData)

	// Calculate and set checksum
	pageBytes := page.ToBytes()
	page.Header.PDChecksum = pgChecksumBlock(pageBytes, blockNumber)

	return page
}

// ToBytes serializes the PostgreSQL page to bytes
func (p *PostgreSQLPage) ToBytes() []byte {
	buf := make([]byte, PostgreSQLPageSize)

	binary.LittleEndian.PutUint64(buf[0:], p.Header.PDLsnXLogRecPtr)
	binary.LittleEndian.PutUint16(buf[8:], p.Header.PDChecksum)
	binary.LittleEndian.PutUint16(buf[10:], p.Header.PDFlags)
	binary.LittleEndian.PutUint16(buf[12:], p.Header.PDLower)
	binary.LittleEndian.PutUint16(buf[14:], p.Header.PDUpper)
	binary.LittleEndian.PutUint16(buf[16:], p.Header.PDSpecial)
	binary.LittleEndian.PutUint16(buf[18:], p.Header.PDPageSizeVer)
	binary.LittleEndian.PutUint32(buf[20:], p.Header.PDPruneXid)

	copy(buf[24:], p.Data)

	return buf
}

// FromBytes deserializes bytes to PostgreSQL page
func (p *PostgreSQLPage) FromBytes(data []byte) {
	p.Header.PDLsnXLogRecPtr = binary.LittleEndian.Uint64(data[0:])
	p.Header.PDChecksum = binary.LittleEndian.Uint16(data[8:])
	p.Header.PDFlags = binary.LittleEndian.Uint16(data[10:])
	p.Header.PDLower = binary.LittleEndian.Uint16(data[12:])
	p.Header.PDUpper = binary.LittleEndian.Uint16(data[14:])
	p.Header.PDSpecial = binary.LittleEndian.Uint16(data[16:])
	p.Header.PDPageSizeVer = binary.LittleEndian.Uint16(data[18:])
	p.Header.PDPruneXid = binary.LittleEndian.Uint32(data[20:])

	p.Data = make([]byte, PostgreSQLPageSize-24)
	copy(p.Data, data[24:])
}

// VerifyChecksum verifies the page checksum
func (p *PostgreSQLPage) VerifyChecksum(blockNumber uint32) bool {
	pageBytes := p.ToBytes()
	// Zero out the checksum field for calculation
	binary.LittleEndian.PutUint16(pageBytes[8:], 0)
	expectedChecksum := pgChecksumBlock(pageBytes, blockNumber)
	return p.Header.PDChecksum == expectedChecksum
}

// ============================================================================
// MySQL/InnoDB Checksum Implementation
// Based on: storage/innobase/ut/ut0crc32.cc
// InnoDB supports multiple checksum algorithms: crc32, innodb, none
// ============================================================================

// innodbChecksumCRC32 calculates InnoDB CRC32 checksum (default in MySQL 5.7+)
func innodbChecksumCRC32(data []byte) uint32 {
	// InnoDB uses CRC32-C (Castagnoli) polynomial
	// Go's standard crc32 with IEEE polynomial is used here for simplicity
	table := crc32.MakeTable(crc32.Castagnoli)

	// Calculate CRC32 of the page body (excluding first 4 bytes and last 8 bytes)
	if len(data) < 12 {
		return 0
	}
	return crc32.Checksum(data[4:len(data)-8], table)
}

// innodbChecksumInnoDB calculates the legacy InnoDB checksum
func innodbChecksumInnoDB(data []byte) uint32 {
	// Legacy InnoDB checksum algorithm (used before MySQL 5.6.3)
	// This is a simple sum-based checksum
	var checksum uint32 = 0

	for i := 0; i < len(data); i += 4 {
		if i+4 <= len(data) {
			val := binary.LittleEndian.Uint32(data[i:])
			checksum += val
		}
	}

	return checksum
}

// NewInnoDBPage creates a new InnoDB page with sample data
func NewInnoDBPage(pageNumber uint32, spaceID uint32, userData string) *InnoDBPage {
	page := &InnoDBPage{
		Header: InnoDBPageHeader{
			PageNumber: pageNumber,
			PrevPage:   0xFFFFFFFF, // No previous page
			NextPage:   0xFFFFFFFF, // No next page
			LSN:        0x0000000100000001,
			PageType:   17855, // FIL_PAGE_INDEX
			FlushLSN:   0,
			SpaceID:    spaceID,
		},
		Data: make([]byte, InnoDBPageSize-38-8), // Header=38, Trailer=8
		Trailer: InnoDBPageTrailer{
			LSNLow: 0x00000001, // Low 4 bytes of LSN
		},
	}

	// Write user data
	copy(page.Data, userData)

	// Calculate and set checksums
	pageBytes := page.ToBytes()
	page.Header.Checksum = innodbChecksumCRC32(pageBytes)
	page.Trailer.OldChecksum = innodbChecksumInnoDB(pageBytes)

	return page
}

// ToBytes serializes the InnoDB page to bytes
func (p *InnoDBPage) ToBytes() []byte {
	buf := make([]byte, InnoDBPageSize)

	// Header (38 bytes)
	binary.LittleEndian.PutUint32(buf[0:], p.Header.Checksum)
	binary.LittleEndian.PutUint32(buf[4:], p.Header.PageNumber)
	binary.LittleEndian.PutUint32(buf[8:], p.Header.PrevPage)
	binary.LittleEndian.PutUint32(buf[12:], p.Header.NextPage)
	binary.LittleEndian.PutUint64(buf[16:], p.Header.LSN)
	binary.LittleEndian.PutUint16(buf[24:], p.Header.PageType)
	binary.LittleEndian.PutUint64(buf[26:], p.Header.FlushLSN)
	binary.LittleEndian.PutUint32(buf[34:], p.Header.SpaceID)

	// Data
	copy(buf[38:], p.Data)

	// Trailer (last 8 bytes)
	binary.LittleEndian.PutUint32(buf[InnoDBPageSize-8:], p.Trailer.OldChecksum)
	binary.LittleEndian.PutUint32(buf[InnoDBPageSize-4:], p.Trailer.LSNLow)

	return buf
}

// FromBytes deserializes bytes to InnoDB page
func (p *InnoDBPage) FromBytes(data []byte) {
	p.Header.Checksum = binary.LittleEndian.Uint32(data[0:])
	p.Header.PageNumber = binary.LittleEndian.Uint32(data[4:])
	p.Header.PrevPage = binary.LittleEndian.Uint32(data[8:])
	p.Header.NextPage = binary.LittleEndian.Uint32(data[12:])
	p.Header.LSN = binary.LittleEndian.Uint64(data[16:])
	p.Header.PageType = binary.LittleEndian.Uint16(data[24:])
	p.Header.FlushLSN = binary.LittleEndian.Uint64(data[26:])
	p.Header.SpaceID = binary.LittleEndian.Uint32(data[34:])

	p.Data = make([]byte, InnoDBPageSize-38-8)
	copy(p.Data, data[38:InnoDBPageSize-8])

	p.Trailer.OldChecksum = binary.LittleEndian.Uint32(data[InnoDBPageSize-8:])
	p.Trailer.LSNLow = binary.LittleEndian.Uint32(data[InnoDBPageSize-4:])
}

// VerifyChecksum verifies the InnoDB page checksum
func (p *InnoDBPage) VerifyChecksum() bool {
	pageBytes := p.ToBytes()
	// Zero out the checksum field for calculation
	binary.LittleEndian.PutUint32(pageBytes[0:], 0)
	expectedChecksum := innodbChecksumCRC32(pageBytes)
	return p.Header.Checksum == expectedChecksum
}

// ============================================================================
// Demo Functions
// ============================================================================

func demonstratePostgreSQLChecksum() {
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println("PostgreSQL Data Page Checksum Demo")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println()

	// Create a page with sample data
	blockNumber := uint32(42)
	userData := "Hello, this is sample row data stored in PostgreSQL!"
	page := NewPostgreSQLPage(blockNumber, userData)

	fmt.Printf("Created PostgreSQL page (block %d)\n", blockNumber)
	fmt.Printf("  Page size: %d bytes\n", PostgreSQLPageSize)
	fmt.Printf("  Checksum: 0x%04X\n", page.Header.PDChecksum)
	fmt.Printf("  User data: %s\n", string(page.Data[:len(userData)]))
	fmt.Println()

	// Verify original checksum
	if page.VerifyChecksum(blockNumber) {
		fmt.Println("[OK] Original page checksum is VALID")
	} else {
		fmt.Println("[ERROR] Original page checksum is INVALID")
	}
	fmt.Println()

	// Simulate corruption - modify one byte in the data area
	fmt.Println("--- Simulating Data Corruption ---")
	originalByte := page.Data[10]
	page.Data[10] = 0xFF // Corrupt a byte
	fmt.Printf("  Modified byte at position 10: 0x%02X -> 0xFF\n", originalByte)
	fmt.Println()

	// Verify corrupted page
	if page.VerifyChecksum(blockNumber) {
		fmt.Println("[ERROR] Corrupted page checksum is VALID (corruption not detected!)")
	} else {
		fmt.Println("[OK] Corrupted page checksum is INVALID (corruption detected!)")
	}

	// Calculate what the new checksum should be
	pageBytes := page.ToBytes()
	binary.LittleEndian.PutUint16(pageBytes[8:], 0)
	newChecksum := pgChecksumBlock(pageBytes, blockNumber)
	fmt.Printf("  Current checksum: 0x%04X\n", page.Header.PDChecksum)
	fmt.Printf("  Expected checksum: 0x%04X\n", newChecksum)
	fmt.Println()

	// Show PostgreSQL error message that would appear
	fmt.Println("PostgreSQL would report:")
	fmt.Printf(`  WARNING:  page verification failed, calculated checksum %d but expected %d
  ERROR:  invalid page in block %d of relation base/12345/16384
`, newChecksum, page.Header.PDChecksum, blockNumber)
	fmt.Println()
}

func demonstrateInnoDBChecksum() {
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println("MySQL/InnoDB Data Page Checksum Demo")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println()

	// Create a page with sample data
	pageNumber := uint32(100)
	spaceID := uint32(1)
	userData := "Hello, this is sample row data stored in InnoDB!"
	page := NewInnoDBPage(pageNumber, spaceID, userData)

	fmt.Printf("Created InnoDB page (page %d, space %d)\n", pageNumber, spaceID)
	fmt.Printf("  Page size: %d bytes\n", InnoDBPageSize)
	fmt.Printf("  CRC32 Checksum (header): 0x%08X\n", page.Header.Checksum)
	fmt.Printf("  Legacy Checksum (trailer): 0x%08X\n", page.Trailer.OldChecksum)
	fmt.Printf("  User data: %s\n", string(page.Data[:len(userData)]))
	fmt.Println()

	// Verify original checksum
	if page.VerifyChecksum() {
		fmt.Println("[OK] Original page checksum is VALID")
	} else {
		fmt.Println("[ERROR] Original page checksum is INVALID")
	}
	fmt.Println()

	// Simulate corruption - modify data in the middle of the page
	fmt.Println("--- Simulating Data Corruption ---")
	originalByte := page.Data[20]
	page.Data[20] = 0xDE // Corrupt a byte
	page.Data[21] = 0xAD // Corrupt another byte
	fmt.Printf("  Modified bytes at position 20-21: 0x%02X... -> 0xDEAD\n", originalByte)
	fmt.Println()

	// Verify corrupted page
	if page.VerifyChecksum() {
		fmt.Println("[ERROR] Corrupted page checksum is VALID (corruption not detected!)")
	} else {
		fmt.Println("[OK] Corrupted page checksum is INVALID (corruption detected!)")
	}

	// Calculate what the new checksum should be
	pageBytes := page.ToBytes()
	binary.LittleEndian.PutUint32(pageBytes[0:], 0)
	newChecksum := innodbChecksumCRC32(pageBytes)
	fmt.Printf("  Current checksum: 0x%08X\n", page.Header.Checksum)
	fmt.Printf("  Expected checksum: 0x%08X\n", newChecksum)
	fmt.Println()

	// Show MySQL error message that would appear
	fmt.Println("MySQL/InnoDB would report:")
	fmt.Printf(`  InnoDB: Page checksum mismatch in space %d page %d
  InnoDB: Stored checksum: 0x%08X, calculated: 0x%08X
  InnoDB: Database page corruption on disk or a failed file read
`, spaceID, pageNumber, page.Header.Checksum, newChecksum)
	fmt.Println()
}

func demonstrateChecksumBypass() {
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println("Demonstrating Checksum Recalculation (Bypass Detection)")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println()

	fmt.Println("IMPORTANT: This demonstrates why checksums alone don't prevent")
	fmt.Println("malicious modification - only unintentional corruption detection.")
	fmt.Println()

	// Create an InnoDB page
	pageNumber := uint32(200)
	spaceID := uint32(2)
	originalData := "Original credit: $100.00"
	page := NewInnoDBPage(pageNumber, spaceID, originalData)

	fmt.Printf("Original page data: %s\n", originalData)
	fmt.Printf("Original checksum: 0x%08X\n", page.Header.Checksum)
	fmt.Println()

	// Maliciously modify the data
	modifiedData := "Modified credit: $999999.99"
	copy(page.Data, modifiedData)

	fmt.Println("--- Malicious Modification ---")
	fmt.Printf("Modified page data: %s\n", modifiedData)

	// Verify - should fail
	if !page.VerifyChecksum() {
		fmt.Println("[DETECTED] Modification detected via checksum mismatch!")
	}
	fmt.Println()

	// Now recalculate the checksum (what a malicious actor would do)
	fmt.Println("--- Recalculating Checksum ---")
	pageBytes := page.ToBytes()
	binary.LittleEndian.PutUint32(pageBytes[0:], 0)
	newChecksum := innodbChecksumCRC32(pageBytes)
	page.Header.Checksum = newChecksum
	page.Trailer.OldChecksum = innodbChecksumInnoDB(page.ToBytes())

	fmt.Printf("New checksum: 0x%08X\n", page.Header.Checksum)

	// Verify - will pass now!
	if page.VerifyChecksum() {
		fmt.Println("[BYPASSED] Page passes checksum verification despite modification!")
	}
	fmt.Println()

	fmt.Println("LESSON: Checksums detect accidental corruption (bit flips, disk errors),")
	fmt.Println("but do NOT prevent intentional tampering. For tamper-detection, use:")
	fmt.Println("  - Cryptographic MACs (HMAC-SHA256)")
	fmt.Println("  - Digital signatures")
	fmt.Println("  - Transparent Data Encryption (TDE)")
	fmt.Println("  - File system integrity monitoring (AIDE, Tripwire)")
	fmt.Println()
}

func printDatabaseChecksumInfo() {
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println("Database Checksum Reference Information")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 70)))
	fmt.Println()

	fmt.Println("PostgreSQL Checksums:")
	fmt.Println("  - Enable: initdb --data-checksums  OR  pg_checksums -e")
	fmt.Println("  - Check status: SHOW data_checksums;")
	fmt.Println("  - Algorithm: FNV-1a based, SIMD optimized")
	fmt.Println("  - Page size: 8KB")
	fmt.Println("  - Checksum location: Bytes 8-9 of page header")
	fmt.Println("  - Verify: pg_verify_checksums (PostgreSQL 11+)")
	fmt.Println()

	fmt.Println("MySQL/InnoDB Checksums:")
	fmt.Println("  - Enable: innodb_checksum_algorithm = crc32 (default)")
	fmt.Println("  - Options: crc32, innodb, none, strict_crc32, strict_innodb, strict_none")
	fmt.Println("  - Page size: 16KB (configurable: 4KB, 8KB, 16KB, 32KB, 64KB)")
	fmt.Println("  - Checksum locations:")
	fmt.Println("      Header bytes 0-3: CRC32 checksum")
	fmt.Println("      Trailer bytes 16376-16379: Old-style checksum")
	fmt.Println("  - Verify: innochecksum utility")
	fmt.Println()

	fmt.Println("Useful Commands:")
	fmt.Println()
	fmt.Println("  PostgreSQL:")
	fmt.Println("    # Check if checksums are enabled")
	fmt.Println("    psql -c \"SHOW data_checksums;\"")
	fmt.Println()
	fmt.Println("    # Verify all checksums (requires server stopped)")
	fmt.Println("    pg_checksums -c -D /var/lib/postgresql/data")
	fmt.Println()
	fmt.Println("  MySQL:")
	fmt.Println("    # Check page checksums")
	fmt.Println("    innochecksum /var/lib/mysql/dbname/tablename.ibd")
	fmt.Println()
	fmt.Println("    # Rewrite checksums")
	fmt.Println("    innochecksum --write=crc32 tablename.ibd")
	fmt.Println()
}

func main() {
	printDatabaseChecksumInfo()
	fmt.Println()

	demonstratePostgreSQLChecksum()
	fmt.Println()

	demonstrateInnoDBChecksum()
	fmt.Println()

	demonstrateChecksumBypass()
}
