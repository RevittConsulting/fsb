package fsb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"reflect"
)

func txBegin(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func txCommit(tx pgx.Tx, ctx context.Context) error {
	err := tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func txDefer(tx pgx.Tx, ctx context.Context) {
	err := tx.Rollback(ctx)
	if err != nil {
		if !errors.Is(err, pgx.ErrTxClosed) {
			_ = fmt.Errorf("error rolling back transaction: %v", err)
		}
	}
}

type ErrResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	if isNil(v) {
		writeErr(w, fmt.Errorf("not found"), http.StatusNotFound)
		return
	}

	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeErr(w http.ResponseWriter, err error, code int) {
	switch err.(type) {
	default:
		w.WriteHeader(code)
		writeJSON(w, ErrResponse{Error: err.Error()})
	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	value := reflect.ValueOf(i)
	kind := value.Kind()
	return (kind == reflect.Ptr || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Func || kind == reflect.Chan || kind == reflect.Interface) && value.IsNil()
}
