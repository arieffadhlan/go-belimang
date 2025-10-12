package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func PrepareStatements(ctx context.Context, pool *pgxpool.Pool) error {
	connection, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	c := connection.Conn()
	log.Info().Msg("Preparing statements...")

	statements := map[string]string{
		// Auth
		"createUser": `
			INSERT INTO users (
        email, 
        username, 
        password, 
        is_admin
      )
			VALUES ($1, $2, $3, $4)
			RETURNING id, is_admin
    `,
		"getUserByMailAddr": `SELECT id, password, is_admin FROM users WHERE email = $1 LIMIT 1`,
		"getUserByUsername": `SELECT id, password, is_admin FROM users WHERE username = $1 LIMIT 1`,

		// Merchant
		"getMerchantById": `SELECT id FROM merchants WHERE id = $1`,
		"createMerchant": `
			INSERT INTO merchants (
        name, 
        imageurl, 
        category, 
        location
      )
			VALUES (
        $1, $2, $3, 
        ST_SetSRID(ST_MakePoint($4, $5), 4326)::GEOGRAPHY
      )
			RETURNING id
		`,
		"createMercItem": `
			INSERT INTO items (merchant_id, name, price, imageurl, category)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`,

		// Purchase
		"getAllMerchantByIDs": `
			SELECT id, name, imageurl, category,
		       ST_X(location::geometry) AS lon,
		       ST_Y(location::geometry) AS lat,
		       created_at
      FROM merchants WHERE id = ANY($1)
    `,
		"getAllMercItemByIDs": `
			SELECT id, merchant_id, name, price, imageurl, category, created_at
      FROM items
      WHERE id = ANY($1)
    `,
		"getEstimateDataByID": `SELECT id, user_id, created_at FROM estimates WHERE id = $1`,
		"createEstimateBatch": `INSERT INTO estimates (user_id) VALUES ($1) RETURNING id`,
		"createOrderFromEsID": `INSERT INTO orders (estimate_id) VALUES ($1)`,
	}

	for name, sql := range statements {
		if _, err := c.Prepare(ctx, name, sql); err != nil {
			log.Error().Err(err).Str("stmt", name).Msg("failed to prepare statement")
			return err
		}
		log.Info().Str("stmt", name).Msg("prepared")
	}

	log.Info().Msg("All statements prepared successfully ✅")
	return nil
}
