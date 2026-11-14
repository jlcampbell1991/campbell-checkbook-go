package db

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    "github.com/google/uuid"
    "github.com/lib/pq"
)

type Video struct {
    ID              uuid.UUID  `db:"id"`
    Name            string     `db:"name"`
    URL             string     `db:"url"`
    Channel         string     `db:"channel"`
    Active          bool       `db:"active"`
    Label           *string    `db:"label"`
    ChannelID       *string    `db:"channel_id"`
    PublishedDate   *time.Time `db:"published_date"`
    DurationSeconds *int       `db:"duration_seconds"`
    CreatedAt       time.Time  `db:"created_at"`
    UpdatedAt       time.Time  `db:"updated_at"`
}

type VideoRepository interface {
    SelectActiveVideos(ids *[]uuid.UUID) ([]Video, error)
    SelectLabels() ([]string, error)
    SelectActiveVideosByLabel(label string, ids *[]uuid.UUID) ([]Video, error)
    SelectByChannelId(channelId string, ids *[]uuid.UUID) ([]Video, error)
    SelectById(id uuid.UUID) (*Video, error)
    Toggle(active bool, id uuid.UUID) error
}

type videoRepository struct {
    database *sql.DB
    limit    int
}

func NewVideoRepo(database *sql.DB, limit int) VideoRepository {
    return videoRepository{database: database, limit: limit}
}

func convertToVids(rows *sql.Rows) ([]Video, error) {
    var videos []Video
    for rows.Next() {
        var item Video
        err := rows.Scan(
            &item.Name,
            &item.URL,
            &item.Channel,
            &item.Active,
            &item.Label,
            &item.CreatedAt,
            &item.UpdatedAt,
            &item.ID,
            &item.ChannelID,
            &item.PublishedDate,
            &item.DurationSeconds,
        )
        if err != nil {
            log.Fatal(err)
        }
        videos = append(videos, item)
    }
    return videos, nil
}

var videoFields = `name,
                url,
                channel,
                active,
                label,
                created_at,
                updated_at,
                id,
                channel_id,
                published_date,
                duration_seconds`

func (repo videoRepository) SelectActiveVideos(ids *[]uuid.UUID) ([]Video, error) {
    fmt.Println("Selecting all active videos")
    baseQuery := `
        SELECT ` + videoFields + `
        FROM videos
        WHERE active = true
    `
    args := []any{}
    argPos := 1
    if ids != nil {
        strs := make([]string, len(*ids))
        for i, id := range *ids {
            strs[i] = id.String()
        }
        baseQuery += fmt.Sprintf(" AND NOT (id = ANY($%d))", argPos)
        args = append(args, pq.Array(strs))
        argPos++
    }
    baseQuery += fmt.Sprintf(" ORDER BY RANDOM() LIMIT %d", repo.limit)
    fmt.Println(baseQuery)
    rows, err := repo.database.Query(baseQuery, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return convertToVids(rows)
}

func (repo videoRepository) SelectLabels() ([]string, error) {
    fmt.Println("Selecting all labels")
    rows, err := repo.database.Query(`
        select distinct label
        from videos
        where label is not null
        order by label
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var labels []string
    for rows.Next() {
        var label string
        err := rows.Scan(&label)
        if err != nil {
            return nil, err
        }
        labels = append(labels, label)
    }
    return labels, nil
}

func (repo videoRepository) SelectActiveVideosByLabel(label string, ids *[]uuid.UUID) ([]Video, error) {
    fmt.Printf("Select all active videos by label: %v\n", label)
    baseQuery := `
        SELECT ` + videoFields + `
        FROM videos
        WHERE active = true
        AND label = $1`
    args := []any{label}
    argPos := 2
    if ids != nil {
        strs := make([]string, len(*ids))
        for i, id := range *ids {
            strs[i] = id.String()
        }
        baseQuery += fmt.Sprintf(" AND NOT (id = ANY($%d))", argPos)
        args = append(args, pq.Array(strs))
        argPos++
    }
    baseQuery += fmt.Sprintf(" ORDER BY RANDOM() LIMIT %d", repo.limit)
    fmt.Println(baseQuery)
    rows, err := repo.database.Query(baseQuery, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return convertToVids(rows)
}

func (repo videoRepository) SelectByChannelId(channelId string, ids *[]uuid.UUID) ([]Video, error) {
    fmt.Printf("Select all active videos by channelId: %v\n", channelId)
    baseQuery := `
        SELECT ` + videoFields + `
        FROM videos
        WHERE channel_id = $1`
    args := []any{channelId}
    argPos := 2
    if ids != nil {
        strs := make([]string, len(*ids))
        for i, id := range *ids {
            strs[i] = id.String()
        }
        baseQuery += fmt.Sprintf(" AND NOT (id = ANY($%d))", argPos)
        args = append(args, pq.Array(strs))
        argPos++
    }
    baseQuery += fmt.Sprintf(" ORDER BY active DESC, published_date DESC LIMIT %d", repo.limit)
    fmt.Println(baseQuery)
    rows, err := repo.database.Query(baseQuery, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return convertToVids(rows)
}

func (repo videoRepository) SelectById(id uuid.UUID) (*Video, error) {
    fmt.Printf("Select video by id: %v\n", id)
    query := `SELECT ` + videoFields + ` FROM videos WHERE id = $1`
    fmt.Println(query)
    row := repo.database.QueryRow(query, id)
    var v Video
    err := row.Scan(
        &v.Name,
        &v.URL,
        &v.Channel,
        &v.Active,
        &v.Label,
        &v.CreatedAt,
        &v.UpdatedAt,
        &v.ID,
        &v.ChannelID,
        &v.PublishedDate,
        &v.DurationSeconds,
    )
    if err != nil {
        return nil, err
    }
    return &v, nil
}

func (repo videoRepository) Toggle(active bool, id uuid.UUID) error {
    fmt.Printf("Setting active = %v for video %v\n", active, id)
    query := `UPDATE videos SET active = $1 WHERE id = $2`
    fmt.Println(query)
    _, err := repo.database.Exec(query, active, id)
    return err
}
