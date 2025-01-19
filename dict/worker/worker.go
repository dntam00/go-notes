package worker

import (
	"database/sql"
	"fmt"
	"play-around/dict/html"
	"play-around/dict/model"
	"play-around/dict/monitor"
	"time"
)

type Worker struct {
	ch            chan model.Entry
	entriesCh     chan []model.Entry
	numbers       int
	db            *sql.DB
	htmlProcessor *html.Processor
}

func New(numbers int, queueSize int, db *sql.DB, htmlProcessor *html.Processor) *Worker {
	return &Worker{
		ch:            make(chan model.Entry, queueSize),
		entriesCh:     make(chan []model.Entry, queueSize),
		numbers:       numbers,
		db:            db,
		htmlProcessor: htmlProcessor,
	}
}

func (w *Worker) AddEntry(entry model.Entry) {
	w.ch <- entry
}

func (w *Worker) AddEntries(entry []model.Entry) {
	w.entriesCh <- entry
}

func (w *Worker) StartWorkerEntries() {
	for i := 0; i < w.numbers; i++ {
		fmt.Println("start worker", i)
		go func() {
			for {
				entries := <-w.entriesCh
				for _, entry := range entries {
					w.exec(entry)
				}
			}
		}()
	}
}

func (w *Worker) StartWorker() {
	for i := 0; i < w.numbers; i++ {
		fmt.Println("start worker", i)
		go func() {
			for {
				entry := <-w.ch
				w.exec(entry)
			}
		}()
	}
}

func (w *Worker) exec(entry model.Entry) {
	query :=
		`
				INSERT INTO word (word, content)
				VALUES (?, ?)
				ON DUPLICATE KEY UPDATE word = (word), content = (content)
	`
	start := time.Now()
	word := entry.Word[:len(entry.Word)-2]
	result, err := w.db.Exec(query, word, entry.Content[:len(entry.Content)-2])
	monitor.RecordCommand("execute", time.Since(start).Milliseconds())
	if err != nil {
		fmt.Println("error inserting entry", err)
		return
	}

	affected, err := result.RowsAffected()

	if err != nil {
		fmt.Println("get affected row failed", err)
		return
	}

	var wordId int64
	if affected == 0 {
		getQuery := `
			SELECT ID FROM word WHERE word = ?
		`
		startQuery := time.Now()
		err := w.db.QueryRow(getQuery, word).Scan(&wordId)
		monitor.RecordCommand("query", time.Since(startQuery).Milliseconds())
		if err != nil {
			fmt.Println("get word on existing failed", err)
			return
		}
	} else {
		wordId, err = result.LastInsertId()
		if err != nil {
			fmt.Println("get last inserted failed", err)
			return
		}
	}
	go w.htmlProcessor.Process(word, wordId, entry.Content)
}
