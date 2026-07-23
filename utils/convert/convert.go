package convert

import "database/sql"

func UintToNullInt64(v uint) any {
	if v == 0 {
		return nil
	}
	return int64(v)
}

func SQLiteNullInt64ToUint(v sql.NullInt64) uint {
	if !v.Valid {
		return 0
	}
	return uint(v.Int64)
}

func IntToNullInt64(v int) any {
	if v == 0 {
		return nil
	}
	return int64(v)
}

func SQLiteNullInt64ToInt(v sql.NullInt64) int {
	if !v.Valid {
		return 0
	}
	return int(v.Int64)
}
