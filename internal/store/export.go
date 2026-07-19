package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Au1rxx/free-vpn-subscriptions/pkg/node"
)

type ExportQuery struct {
	Since                        time.Time
	Grades, Protocols, Countries []string
	Limit, Offset                int
}

type ExportMeta struct {
	ConfigID          uint64
	Grade             string
	Score             int
	Country           string
	NetworkClass      string
	ConfigFingerprint [32]byte
	Reason            string
}

func ListExportable(ctx context.Context, db *sql.DB, query ExportQuery) ([]*node.Node, []ExportMeta, error) {
	if query.Limit <= 0 {
		query.Limit = 2000
	}
	if query.Limit > 10000 || query.Offset < 0 {
		return nil, nil, fmt.Errorf("invalid export pagination")
	}
	var where []string
	args := []any{}
	where = append(where, "n.is_exportable=TRUE", "s.availability_state IN ('available','degraded')")
	if !query.Since.IsZero() {
		where = append(where, "n.last_success_at>=?")
		args = append(args, query.Since.UTC())
	}
	if len(query.Grades) > 0 {
		where = append(where, "s.quality_grade IN ("+scalarPlaceholders(len(query.Grades))+")")
		for _, v := range query.Grades {
			args = append(args, v)
		}
	}
	if len(query.Protocols) > 0 {
		where = append(where, "n.protocol IN ("+scalarPlaceholders(len(query.Protocols))+")")
		for _, v := range query.Protocols {
			args = append(args, v)
		}
	}
	if len(query.Countries) > 0 {
		where = append(where, "c.exit_country IN ("+scalarPlaceholders(len(query.Countries))+")")
		for _, v := range query.Countries {
			args = append(args, v)
		}
	}
	args = append(args, query.Limit, query.Offset)
	rows, err := db.QueryContext(ctx, `SELECT n.node_config_id, n.normalized_config, n.config_fingerprint,
		s.quality_grade, s.quality_score, COALESCE(s.exit_country,''), COALESCE(c.provider_class,''),
		COALESCE(s.latency_p50_ms,0)
		FROM node_configs n JOIN node_current_status s ON s.node_config_id=n.node_config_id
		LEFT JOIN node_classifications c ON c.node_config_id=n.node_config_id WHERE `+strings.Join(where, " AND ")+`
		ORDER BY FIELD(s.quality_grade,'S','A','B','C','D','U'), s.quality_score DESC,
		n.protocol, COALESCE(s.latency_p50_ms,4294967295), n.config_fingerprint LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var nodes []*node.Node
	var metadata []ExportMeta
	for rows.Next() {
		var body, fingerprint []byte
		var meta ExportMeta
		var country string
		var latency int
		if err := rows.Scan(&meta.ConfigID, &body, &fingerprint, &meta.Grade, &meta.Score, &country, &meta.NetworkClass, &latency); err != nil {
			return nil, nil, err
		}
		var configured node.Node
		if err := json.Unmarshal(body, &configured); err != nil {
			return nil, nil, fmt.Errorf("decode node %d: %w", meta.ConfigID, err)
		}
		configured.Country, configured.LatencyMS = country, latency
		meta.Country = country
		configured.Name = fmt.Sprintf("%s-%d", configured.Protocol, meta.ConfigID)
		if len(fingerprint) != 32 {
			return nil, nil, fmt.Errorf("node %d has invalid fingerprint", meta.ConfigID)
		}
		copy(meta.ConfigFingerprint[:], fingerprint)
		meta.Reason = "verified_" + strings.ToLower(meta.Grade)
		nodes = append(nodes, &configured)
		metadata = append(metadata, meta)
	}
	return nodes, metadata, rows.Err()
}
