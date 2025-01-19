package html

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"play-around/dict/model"
	"play-around/dict/monitor"
	"strings"
	"time"
)

const (
	entrySelector         = ".ldoceEntry.Entry"
	tailSelector          = ".Tail .Crossref .crossRef"
	bussSelector          = ".bussdictEntry.Entry"
	pronunciationSelector = ".PRON"
	ameFileSelector       = ".amefile"
	breFileSelector       = ".brefile"
	href                  = "href"
	dictRef               = "@@@LINK="
	delimiter             = "\r\n"
	topicSelector         = ".browse_results"
	derivedSelector       = ".DERIV"
	crossRefSelector      = ".Crossref .crossRef"
	relatedRefSelector    = ".RELATEDWD"
	centurySelector       = ".CENTURY"
	refHWDSelector        = ".REFHWD"
)

type Processor struct {
	database *sql.DB
}

func New(database *sql.DB) *Processor {
	return &Processor{
		database: database,
	}
}

func (h *Processor) Process(word string, wordId int64, content string) {
	proceeded := false

	// 1. pure reference case
	if strings.HasPrefix(content, dictRef) {
		monitor.RecordCommand("link", 0)
		proceeded = true
		refWord := strings.TrimSuffix(content, delimiter)
		refWord = strings.TrimPrefix(refWord, dictRef)
		h.updateRef(word, refWord)
		return
	}

	reader := strings.NewReader(content)
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("create document failed", err)
	}

	// 2. reference to pure form of word
	refEntrySelection := document.Find(".ref_entry .ref_item a")
	if len(refEntrySelection.Nodes) > 0 {
		proceeded = true
		var refWord string
		refEntrySelection.Each(func(i int, s *goquery.Selection) {
			h.validateRef(word, s.Text())
			refWord = fmt.Sprintf("%s/%s", refWord, s.Text())
		})
		refWord = strings.TrimSpace(refWord)
		refWord = strings.TrimSuffix(refWord, delimiter)
		if len(refWord) > 0 {
			h.updateRef(word, refWord)
			return
		}
	}

	// 3. reference to tail form of word
	refEntrySelection = document.Find(tailSelector)
	if len(refEntrySelection.Nodes) > 0 {
		proceeded = true
		var refWord string
		refEntrySelection.Each(func(i int, s *goquery.Selection) {
			val, _ := s.Attr("href")
			refWord = strings.TrimPrefix(val, "entry://")
		})
		h.updateRef(word, refWord)
	}

	find := document.Find(entrySelector)
	if len(find.Nodes) != 0 {
		proceeded = true
		find.Each(func(i int, s *goquery.Selection) {
			//fmt.Printf("word type:%s\n", s.Find(".POS").Nodes[0].FirstChild.Data)
			wordFamily := ""
			familyNode := s.Find(".POS")
			if len(familyNode.Nodes) > 0 {
				wordFamily = familyNode.Nodes[0].FirstChild.Data
			}

			pronunciation := ""
			s.Find(pronunciationSelector).Each(func(i int, s *goquery.Selection) {
				pronunciation = fmt.Sprintf("%s/%s/", pronunciation, s.Text())
				//fmt.Printf("pronunciation %d: %s\n", i, s.Text())
			})

			breFile := ""
			s.Find(ameFileSelector).Each(func(i int, s *goquery.Selection) {
				attrs := s.Nodes[0].Attr
				for _, attr := range attrs {
					if attr.Key == href {
						breFile = attr.Val
						//fmt.Printf("ame sound file %d: %s\n", i, attr.Val)
					}
				}
			})

			ameFile := ""
			s.Find(breFileSelector).Each(func(i int, s *goquery.Selection) {
				attrs := s.Nodes[0].Attr
				for _, attr := range attrs {
					if attr.Key == href {
						ameFile = attr.Val
						//fmt.Printf("bre sound file %d: %s\n", i, attr.Val)
					}
				}
			})

			resource := model.Resource{
				Audio: []model.Audio{
					{Path: ameFile, Language: "en-us"},
					{Path: breFile, Language: "en-gb"},
				},
			}

			resourceData, err := json.Marshal(resource)
			if err != nil {
				log.Fatal("convert resource to json failed", err)
			}

			start := time.Now()

			query := `
			insert into word_family (word_id, word_family, pronunciation, resource)
			values (?, ?, ?, ?)
			`

			wordFamilyResult, err := h.database.Exec(query, wordId, wordFamily, pronunciation, string(resourceData))
			wordFamlilyId, err := wordFamilyResult.LastInsertId()
			if err != nil {
				log.Fatal("insert word family failed: ", wordId, pronunciation, err)
				return
			}

			monitor.RecordCommand("execute_word_family", time.Since(start).Milliseconds())

			if err != nil {
				log.Fatal("insert word family failed: ", wordId, pronunciation, err)
				return
			}
			index := 0
			idSelector := fmt.Sprintf("[id^=\"%s__\"]", word)
			s.Find(idSelector).Each(func(i int, s *goquery.Selection) {
				index++
				defHtml, _ := s.Html()
				var defStr string
				defSelection := s.Find(".DEF")
				if len(defSelection.Nodes) > 0 {
					defSelection.Each(func(i int, a *goquery.Selection) {
						defStr = strings.TrimSpace(a.Text())
						query := `
							INSERT INTO meaning (word_id, word_family_id, meaning, meaning_html, position, meaning_type)
							VALUES (?, ?, ?, ?, ?, ?)
						`
						_, err = h.database.Exec(query, wordId, wordFamlilyId, defStr, defHtml, index, "raw")
						if err != nil {
							log.Fatalf("insert meaning %s of word %s failed: %s", defStr, word, err)
							return
						}
					})
					return
				}

				defSelection = s.Find(derivedSelector)
				if len(defSelection.Nodes) > 0 {
					//
					return
				}

				defSelection = s.Find(crossRefSelector)
				if len(defSelection.Nodes) > 0 {
					defSelection.Each(func(i int, a *goquery.Selection) {
						val, _ := a.Attr("href")
						defStr = strings.TrimPrefix(val, "entry://")
						query := `
							INSERT INTO meaning (word_id, word_family_id, meaning, meaning_html, position, meaning_type)
							VALUES (?, ?, ?, ?, ?, ?)
						`
						_, err = h.database.Exec(query, wordId, wordFamlilyId, defStr, defHtml, index, "cross_ref")
						if err != nil {
							log.Fatalf("insert meaning %s of word %s failed: %s", defStr, word, err)
							return
						}
					})
					return
				}

				defSelection = s.Find(relatedRefSelector)
				if len(defSelection.Nodes) > 0 {
					defSelection.Each(func(i int, a *goquery.Selection) {
						defStr = strings.TrimSpace(a.Text())
						query := `
							INSERT INTO meaning (word_id, word_family_id, meaning, meaning_html, position, meaning_type)
							VALUES (?, ?, ?, ?, ?, ?)
						`
						_, err = h.database.Exec(query, wordId, wordFamlilyId, defStr, defHtml, index, "related_ref")
						if err != nil {
							log.Fatalf("insert meaning %s of word %s failed: %s", defStr, word, err)
							return
						}
					})
					return
				}

				defSelection = s.Find(centurySelector)
				if len(defSelection.Nodes) > 0 {
					return
				}

				defSelection = s.Find(refHWDSelector)
				if len(defSelection.Nodes) > 0 {
					defSelection.Each(func(i int, a *goquery.Selection) {
						val, _ := a.Attr("href")
						defStr = strings.TrimPrefix(val, "entry://")
						query := `
							INSERT INTO meaning (word_id, word_family_id, meaning, meaning_html, position, meaning_type)
							VALUES (?, ?, ?, ?, ?, ?)
						`
						_, err = h.database.Exec(query, wordId, wordFamlilyId, defStr, defHtml, index, "cross_ref")
						if err != nil {
							log.Fatalf("insert meaning %s of word %s failed: %s", defStr, word, err)
							return
						}
					})
					return
				}

				fmt.Println("def not found", word)
			})
		})
	}

	find = document.Find(bussSelector)
	if len(find.Nodes) != 0 {
		proceeded = true
		find.Each(func(i int, s *goquery.Selection) {
			//fmt.Printf("word type:%s\n", s.Find(".POS").Nodes[0].FirstChild.Data)
			wordFamily := ""
			familyNode := s.Find(".POS")
			if len(familyNode.Nodes) > 0 {
				wordFamily = familyNode.Nodes[0].FirstChild.Data
			}

			pronunciation := ""
			s.Find(pronunciationSelector).Each(func(i int, s *goquery.Selection) {
				pronunciation = fmt.Sprintf("%s/%s/", pronunciation, s.Text())
				//fmt.Printf("pronunciation %d: %s\n", i, s.Text())
			})

			query := `
			insert into word_family (word_id, word_family, pronunciation, resource)
			values (?, ?, ?, ?)
			`

			_, err = h.database.Exec(query, wordId, wordFamily, pronunciation, nil)

			if err != nil {
				log.Fatal("insert buss word family failed: ", wordId, pronunciation, err)
				return
			}
		})
	}

	topicFind := document.Find(topicSelector)
	if len(topicFind.Nodes) != 0 {
		proceeded = true
		var refWord string
		topicFind.Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				h.validateRef(word, s.Text())
				refWord = fmt.Sprintf("%s/%s", refWord, s.Text())
			})
		})
		h.updateRef(word, refWord)
	}

	if !proceeded {
		log.Println("has not been processed for", wordId)
		monitor.RecordCommand("execute_notfound", 0)
	}
}

func (h *Processor) validateRef(word, ref string) {
	var refId int64
	err := h.database.QueryRow("SELECT id FROM word w WHERE w.word = ?", ref).Scan(&refId)
	if err != nil {
		log.Printf("ref: %s for word: %s not found", ref, word)
	}
}

func (h *Processor) updateRef(word string, ref string) {
	_, err := h.database.Exec("UPDATE word SET ref_word = ? WHERE word = ?", ref, word)
	if err != nil {
		log.Fatalf("update ref: %s failed for word %s, error: %s", ref, word, err)
	}
}
