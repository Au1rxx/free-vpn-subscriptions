package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ExportMember struct {
	ConfigID   uint64
	Collection string
	Rank       int
	Score      int
	Grade      string
	Reason     string
}

func StartExportRun(ctx context.Context, db *sql.DB, runUUID string, startedAt time.Time) (uint64, error) {
	result, err := db.ExecContext(ctx, `INSERT INTO export_runs
		(run_uuid, rules_version, started_at, export_state) VALUES (?, 'db-export-v1', ?, 'running')`, runUUID, startedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

func FailExportRun(ctx context.Context, db *sql.DB, runID uint64, finishedAt time.Time, cause error) error {
	message := "unknown export failure"
	if cause != nil {
		message = cause.Error()
	}
	if len(message) > 1024 {
		message = message[:1024]
	}
	_, err := db.ExecContext(ctx, `UPDATE export_runs SET finished_at=?, export_state='failed', error_summary=?
		WHERE export_run_id=?`, finishedAt, message, runID)
	return err
}

func CompleteExportRun(ctx context.Context, db *sql.DB, runID uint64, finishedAt time.Time,
	candidates, selected, files int, outputBytes int64, summary any, members []ExportMember) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for start := 0; start < len(members); start += 500 {
		end := start + 500
		if end > len(members) {
			end = len(members)
		}
		values := make([]string, 0, end-start)
		args := make([]any, 0, (end-start)*7)
		for _, member := range members[start:end] {
			values = append(values, "(?,?,?,?,?,?,?)")
			args = append(args, runID, member.ConfigID, member.Collection, member.Rank, member.Score, member.Grade, member.Reason)
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO export_members
			(export_run_id, node_config_id, collection_name, rank_number, quality_score, quality_grade, selection_reason)
			VALUES `+strings.Join(values, ","), args...)
		if err != nil {
			return fmt.Errorf("insert export members: %w", err)
		}
	}
	body, err := json.Marshal(summary)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE export_runs SET finished_at=?, candidate_count=?, selected_count=?,
		file_count=?, output_bytes=?, export_state='complete', summary=?, error_summary=NULL WHERE export_run_id=?`,
		finishedAt, candidates, selected, files, outputBytes, body, runID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
