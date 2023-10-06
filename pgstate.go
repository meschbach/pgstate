package pgstate

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

func EnsureDatabase(ctx context.Context, pgConnection *pgx.ConnConfig, databaseName, password string) error {
	cluster, err := pgx.ConnectConfig(ctx, pgConnection)
	if err != nil {
		return err
	}
	defer cluster.Close(ctx)

	if err := ensureRole(ctx, cluster, databaseName, password); err != nil {
		return errors.Join(errors.New("failed to create role"), err)
	}

	escapedDatabaseName := quoteString(databaseName)
	_, err = countRows(ctx, cluster, "CREATE DATABASE "+escapedDatabaseName+" WITH OWNER "+escapedDatabaseName+" ENCODING 'UTF-8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8'")
	if err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			if strings.HasSuffix(e.Message, "already exists") {
				return nil
			}
		}
		return errors.Join(errors.New("failed to create database"), err)
	}

	return nil
}

func ensureRole(ctx context.Context, connection *pgx.Conn, roleName, secret string) error {
	quotedRole := quoteString(roleName)
	quotedSecret := quoteIdentifier(secret)
	_, err := countRows(ctx, connection, "CREATE ROLE "+quotedRole+" WITH LOGIN PASSWORD "+quotedSecret)
	if err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			if strings.HasSuffix(e.Message, "already exists") {
				if _, err := countRows(ctx, connection, "ALTER ROLE "+quotedRole+" PASSWORD "+quotedSecret); err != nil {
					return err
				}
				return nil
			}
		}
		return errors.Join(errors.New("failed to CREATE ROLE ACL"), err)
	}
	//todo: return created
	return nil
}

func countRows(ctx context.Context, connection *pgx.Conn, query string, args ...any) (count int, err error) {
	result, err := connection.Query(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	defer result.Close()

	count = 0
	for result.Next() {
		count++
	}
	return count, result.Err()
}

// todo: improve
func quoteIdentifier(hostile string) string {
	return "'" + hostile + "'"
}

// todo: improve
func quoteString(hostile string) string {
	return "\"" + hostile + "\""
}

func DestroyDatabase(ctx context.Context, pgConnection *pgx.ConnConfig, databaseName string) error {
	cluster, err := pgx.ConnectConfig(ctx, pgConnection)
	if err != nil {
		return err
	}
	defer cluster.Close(ctx)

	_, err = countRows(ctx, cluster, "DROP DATABASE "+quoteString(databaseName)+" WITH (FORCE)")
	if err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			if strings.HasSuffix(e.Message, "does not exist") {
				return nil
			}
		}
		return err
	}
	return nil
}

func DestroyRole(ctx context.Context, pgConnection *pgx.ConnConfig, roleName string) error {
	cluster, err := pgx.ConnectConfig(ctx, pgConnection)
	if err != nil {
		return err
	}
	defer cluster.Close(ctx)

	_, err = countRows(ctx, cluster, "DROP ROLE "+quoteString(roleName))
	if err != nil {
		if e, ok := err.(*pgconn.PgError); ok {
			if strings.HasSuffix(e.Message, "does not exist") {
				return nil
			}
		}
		return err
	}
	return nil
}
