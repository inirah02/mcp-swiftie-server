package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type PrestoClient struct {
	// Mock in-memory database
	albums []Album
	songs  []Song
	tours  []Tour
}

func NewPrestoClient() *PrestoClient {
	return &PrestoClient{
		albums: getSwiftAlbums(),
		songs:  getSwiftSongs(),
		tours:  getSwiftTours(),
	}
}

func (p *PrestoClient) Query(ctx context.Context, sql string) (*QueryResult, error) {
	start := time.Now()

	// Simulate network latency
	time.Sleep(50 * time.Millisecond)

	// Simple SQL parser (mock)
	sql = strings.ToLower(strings.TrimSpace(sql))

	var result *QueryResult
	var err error

	switch {
	case strings.Contains(sql, "show tables"):
		result = p.showTables()
	case strings.Contains(sql, "albums"):
		result = p.queryAlbums(ctx, sql)
	case strings.Contains(sql, "songs"):
		result = p.querySongs(ctx, sql)
	case strings.Contains(sql, "tours"):
		result = p.queryTours(ctx, sql)
	default:
		err = fmt.Errorf("unsupported query: %s", sql)
	}

	if result != nil {
		result.QueryTime = time.Since(start)
	}

	return result, err
}

func (p *PrestoClient) StreamQuery(ctx context.Context, sql string, batchSize int) (<-chan [][]interface{}, <-chan error) {
	rowsChan := make(chan [][]interface{}, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(rowsChan)
		defer close(errChan)

		result, err := p.Query(ctx, sql)
		if err != nil {
			errChan <- err
			return
		}

		// Stream in batches
		for i := 0; i < len(result.Rows); i += batchSize {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				end := i + batchSize
				if end > len(result.Rows) {
					end = len(result.Rows)
				}

				batch := result.Rows[i:end]
				rowsChan <- batch

				// Simulate streaming delay
				time.Sleep(20 * time.Millisecond)
			}
		}
	}()

	return rowsChan, errChan
}

func (p *PrestoClient) showTables() *QueryResult {
	return &QueryResult{
		Columns: []string{"table_name"},
		Rows: [][]interface{}{
			{"albums"},
			{"songs"},
			{"tours"},
		},
		RowCount: 3,
	}
}

func (p *PrestoClient) queryAlbums(ctx context.Context, sql string) *QueryResult {
	rows := make([][]interface{}, 0, len(p.albums))

	for _, album := range p.albums {
		select {
		case <-ctx.Done():
			return nil
		default:
			rows = append(rows, []interface{}{
				album.ID,
				album.Title,
				album.ReleaseYear,
				album.Era,
				album.Sales,
				album.Genre,
			})
		}
	}

	return &QueryResult{
		Columns:  []string{"id", "title", "release_year", "era", "sales_millions", "genre"},
		Rows:     rows,
		RowCount: len(rows),
	}
}

func (p *PrestoClient) querySongs(ctx context.Context, sql string) *QueryResult {
	rows := make([][]interface{}, 0, len(p.songs))

	for _, song := range p.songs {
		select {
		case <-ctx.Done():
			return nil
		default:
			rows = append(rows, []interface{}{
				song.ID,
				song.AlbumID,
				song.Title,
				song.Duration,
				song.Streams,
				song.ChartPeak,
				song.GrammyNoms,
			})
		}
	}

	return &QueryResult{
		Columns:  []string{"id", "album_id", "title", "duration_seconds", "streams_millions", "chart_peak", "grammy_nominations"},
		Rows:     rows,
		RowCount: len(rows),
	}
}

func (p *PrestoClient) queryTours(ctx context.Context, sql string) *QueryResult {
	rows := make([][]interface{}, 0, len(p.tours))

	for _, tour := range p.tours {
		select {
		case <-ctx.Done():
			return nil
		default:
			rows = append(rows, []interface{}{
				tour.ID,
				tour.Name,
				tour.Year,
				tour.Shows,
				tour.Attendance,
				tour.Revenue,
			})
		}
	}

	return &QueryResult{
		Columns:  []string{"id", "name", "year", "shows", "attendance", "revenue_millions"},
		Rows:     rows,
		RowCount: len(rows),
	}
}

// Mock Data
func getSwiftAlbums() []Album {
	return []Album{
		{"ALB001", "Taylor Swift", 2006, "Country", 5, "Country"},
		{"ALB002", "Fearless", 2008, "Country", 12, "Country Pop"},
		{"ALB003", "Speak Now", 2010, "Country Pop", 6, "Country Pop"},
		{"ALB004", "Red", 2012, "Country Pop", 7, "Pop Rock"},
		{"ALB005", "1989", 2014, "Pop", 10, "Synth Pop"},
		{"ALB006", "Reputation", 2017, "Pop", 4, "Electropop"},
		{"ALB007", "Lover", 2019, "Pop", 3, "Pop"},
		{"ALB008", "Folklore", 2020, "Indie Folk", 3, "Indie Folk"},
		{"ALB009", "Evermore", 2020, "Indie Folk", 2, "Alternative"},
		{"ALB010", "Midnights", 2022, "Synth Pop", 6, "Synth Pop"},
		{"ALB011", "The Tortured Poets Department", 2024, "Alternative", 4, "Alternative Pop"},
	}
}

func getSwiftSongs() []Song {
	return []Song{
		{"SONG001", "ALB002", "Love Story", 236, 1800, 4, 0},
		{"SONG002", "ALB002", "You Belong With Me", 232, 1500, 2, 1},
		{"SONG003", "ALB004", "We Are Never Getting Back Together", 193, 1200, 1, 0},
		{"SONG004", "ALB004", "I Knew You Were Trouble", 219, 1400, 2, 1},
		{"SONG005", "ALB005", "Shake It Off", 219, 3200, 1, 3},
		{"SONG006", "ALB005", "Blank Space", 231, 3000, 1, 2},
		{"SONG007", "ALB005", "Style", 231, 1100, 6, 0},
		{"SONG008", "ALB006", "Look What You Made Me Do", 211, 1600, 1, 0},
		{"SONG009", "ALB007", "ME!", 193, 900, 2, 0},
		{"SONG010", "ALB008", "Cardigan", 239, 800, 1, 1},
		{"SONG011", "ALB008", "Exile", 284, 700, 6, 1},
		{"SONG012", "ALB009", "Willow", 214, 600, 1, 0},
		{"SONG013", "ALB010", "Anti-Hero", 200, 2100, 1, 6},
		{"SONG014", "ALB010", "Lavender Haze", 202, 900, 2, 0},
		{"SONG015", "ALB011", "Fortnight", 228, 1100, 1, 0},
		// Add more for streaming demo
		{"SONG016", "ALB005", "Bad Blood", 211, 800, 1, 1},
		{"SONG017", "ALB005", "Wildest Dreams", 220, 1300, 5, 0},
		{"SONG018", "ALB006", "Delicate", 232, 750, 12, 0},
		{"SONG019", "ALB007", "Lover", 221, 850, 10, 0},
		{"SONG020", "ALB008", "The 1", 210, 500, 27, 0},
	}
}

func getSwiftTours() []Tour {
	return []Tour{
		{"TOUR001", "Fearless Tour", 2009, 118, 1200000, 63.5},
		{"TOUR002", "Speak Now World Tour", 2011, 111, 1600000, 123.0},
		{"TOUR003", "The Red Tour", 2013, 86, 1700000, 150.2},
		{"TOUR004", "The 1989 World Tour", 2015, 85, 2278647, 250.7},
		{"TOUR005", "Reputation Stadium Tour", 2018, 53, 2888892, 345.7},
		{"TOUR006", "The Eras Tour", 2023, 152, 10000000, 2000.0},
	}
}
