package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	s := &Storage{db: db}

	if err := s.init(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return s, nil
}

func (s *Storage) init() error {
	walletsQuery := `
	CREATE TABLE IF NOT EXISTS user_wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		wallet TEXT NOT NULL,
		UNIQUE(chat_id, wallet)
	);`

	marketsQuery := `
	CREATE TABLE IF NOT EXISTS user_markets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		market_id TEXT NOT NULL,
		UNIQUE(chat_id, market_id)
	);`

	if _, err := s.db.Exec(walletsQuery); err != nil {
		return err
	}

	if _, err := s.db.Exec(marketsQuery); err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddWallet(chatID int64, wallet string) error {
	query := `INSERT OR IGNORE INTO user_wallets (chat_id, wallet) VALUES (?, ?)`
	_, err := s.db.Exec(query, chatID, wallet)
	if err != nil {
		return fmt.Errorf("failed to add wallet: %w", err)
	}
	return nil
}

func (s *Storage) RemoveWallet(chatID int64, wallet string) error {
	query := `DELETE FROM user_wallets WHERE chat_id = ? AND wallet = ?`
	result, err := s.db.Exec(query, chatID, wallet)
	if err != nil {
		return fmt.Errorf("failed to remove wallet: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("wallet not found")
	}
	return nil
}

func (s *Storage) GetWallets(chatID int64) ([]string, error) {
	query := `SELECT wallet FROM user_wallets WHERE chat_id = ? ORDER BY wallet`
	rows, err := s.db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query wallets: %w", err)
	}
	defer rows.Close()

	var wallets []string
	for rows.Next() {
		var wallet string
		if err := rows.Scan(&wallet); err != nil {
			return nil, fmt.Errorf("failed to scan wallet: %w", err)
		}
		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return wallets, nil
}

func (s *Storage) GetAllSubs() (map[int64][]string, error) {
	query := `SELECT chat_id, wallet FROM user_wallets ORDER BY chat_id, wallet`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all subscriptions: %w", err)
	}
	defer rows.Close()

	subs := make(map[int64][]string)
	for rows.Next() {
		var chatID int64
		var wallet string
		if err := rows.Scan(&chatID, &wallet); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs[chatID] = append(subs[chatID], wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return subs, nil
}

func (s *Storage) AddMarket(chatID int64, marketID string) error {
	query := `INSERT OR IGNORE INTO user_markets (chat_id, market_id) VALUES (?, ?)`
	_, err := s.db.Exec(query, chatID, marketID)
	if err != nil {
		return fmt.Errorf("failed to add market: %w", err)
	}
	return nil
}

func (s *Storage) RemoveMarket(chatID int64, marketID string) error {
	query := `DELETE FROM user_markets WHERE chat_id = ? AND market_id = ?`
	result, err := s.db.Exec(query, chatID, marketID)
	if err != nil {
		return fmt.Errorf("failed to remove market: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("market not found")
	}
	return nil
}

func (s *Storage) GetMarkets(chatID int64) ([]string, error) {
	query := `SELECT market_id FROM user_markets WHERE chat_id = ? ORDER BY market_id`
	rows, err := s.db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query markets: %w", err)
	}
	defer rows.Close()

	var markets []string
	for rows.Next() {
		var marketID string
		if err := rows.Scan(&marketID); err != nil {
			return nil, fmt.Errorf("failed to scan market: %w", err)
		}
		markets = append(markets, marketID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return markets, nil
}

func (s *Storage) GetAllMarkets() (map[int64][]string, error) {
	query := `SELECT chat_id, market_id FROM user_markets ORDER BY chat_id, market_id`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all markets: %w", err)
	}
	defer rows.Close()

	markets := make(map[int64][]string)
	for rows.Next() {
		var chatID int64
		var marketID string
		if err := rows.Scan(&chatID, &marketID); err != nil {
			return nil, fmt.Errorf("failed to scan market: %w", err)
		}
		markets[chatID] = append(markets[chatID], marketID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return markets, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
