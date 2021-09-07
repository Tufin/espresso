package bq

import (
	"context"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/env"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	EnvKeyBQToken = "BIGQUERY_KEY"
)

type Client interface {
	QueryIterator(q string, params []bigquery.QueryParameter) (Iterator, error)
}

type ClientImpl struct {
	bqClient *bigquery.Client
}

func NewClient(gcpProjectID string) Client {

	if key := env.GetSensitive(EnvKeyBQToken); key != "" {
		conf, err := google.JWTConfigFromJSON([]byte(key), bigquery.Scope)
		if err != nil {
			log.Fatalf("failed to config big-query JWT with %v", err)
		}

		ctx := context.Background()
		client, err := bigquery.NewClient(ctx, gcpProjectID, option.WithTokenSource(conf.TokenSource(ctx)))
		if err != nil {
			log.Fatalf("failed to create bigquery client with %v", err)
		}

		return &ClientImpl{bqClient: client}
	}

	client, err := bigquery.NewClient(context.Background(), gcpProjectID)
	if err != nil {
		log.Fatalf("failed to create bigquery client without token with %v", err)
	}

	return &ClientImpl{bqClient: client}
}

func (client *ClientImpl) QueryIterator(q string, params []bigquery.QueryParameter) (Iterator, error) {

	query := client.bqClient.Query(q)
	query.Parameters = params
	rowIterator, err := query.Read(context.Background())
	if err != nil {
		return nil, err
	}

	return newRawIterator(rowIterator), err
}

func (client *ClientImpl) Query(q string) error {

	_, err := client.bqClient.Query(q).Read(context.Background())

	return err
}

func (client *ClientImpl) GetQueryStats(q string, params []bigquery.QueryParameter) (*bigquery.JobStatistics, error) {

	query := client.bqClient.Query(q)
	query.Parameters = params
	query.QueryConfig.DryRun = true

	job, err := query.Run(context.Background())
	if err != nil {
		log.Errorf("get query stats failed with %v", err)
		return nil, err
	}

	return job.LastStatus().Statistics, nil
}
