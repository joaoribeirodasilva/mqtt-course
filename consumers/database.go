package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Database struct {
	conf     *Configuration
	Client   *mongo.Client
	Database *mongo.Database
}

func NewDatabase(conf *Configuration) *Database {

	d := &Database{}

	d.conf = conf

	return d
}

func (d *Database) Connect() error {

	// create the options
	clientOpts := options.Client()

	// set the server API options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts.SetServerAPIOptions(serverAPI)

	// user credentials (if set)
	if d.conf.Mongo.Username != "" && d.conf.Mongo.Password != "" {
		credential := options.Credential{
			Username:    d.conf.Mongo.Username,
			Password:    d.conf.Mongo.Password,
			PasswordSet: true,
			AuthSource:  "admin",
		}
		clientOpts.SetAuth(credential)
	}

	// connection pool
	clientOpts.SetMinPoolSize(uint64(d.conf.Mongo.MinPoolSize))
	clientOpts.SetMaxPoolSize(uint64(d.conf.Mongo.MaxPoolSize))
	clientOpts.SetReplicaSet(d.conf.Mongo.ReplicaSet)

	// timeouts
	clientOpts.SetTimeout(time.Duration(d.conf.Mongo.TimeoutMS) * time.Millisecond)
	clientOpts.SetConnectTimeout(time.Duration(d.conf.Mongo.ConnectTimeoutMS) * time.Millisecond)
	clientOpts.SetMaxConnIdleTime(time.Duration(d.conf.Mongo.MaxIdleTimeMS) * time.Millisecond)
	clientOpts.SetSocketTimeout(time.Duration(d.conf.Mongo.SocketTimeoutMS) * time.Millisecond)
	clientOpts.SetServerSelectionTimeout(time.Duration(d.conf.Mongo.ServerSelectionTimeoutMS) * time.Millisecond)

	// heartbeat
	clientOpts.SetHeartbeatInterval(time.Duration(d.conf.Mongo.HeartbeatFrequencyMS) * time.Millisecond)

	// write concern
	writeConcern := writeconcern.WriteConcern{
		W:        d.conf.Mongo.WriteConcern.W,
		WTimeout: time.Duration(d.conf.Mongo.WriteConcern.WTimeoutMS) * time.Millisecond,
		Journal:  &d.conf.Mongo.WriteConcern.Journal,
	}
	clientOpts.SetWriteConcern(&writeConcern)

	// read preference
	readPref := readpref.Primary()
	if d.conf.Mongo.ReadPreference.ReadPreference != "" {
		switch strings.ToLower(d.conf.Mongo.ReadPreference.ReadPreference) {
		case "primary":
			// it's already set by default
		case "primarypreferred":
			readPref = readpref.PrimaryPreferred()
		default:
			return fmt.Errorf("ERROR: [DATABASE] invalid read preference %s", d.conf.Mongo.ReadPreference.ReadPreference)
		}
	}
	clientOpts.SetReadPreference(readPref)

	// direct connection
	clientOpts.SetDirect(d.conf.Mongo.DirectConnection)

	// data compression
	compressors := make([]string, 0)
	if d.conf.Mongo.Compressors.Snappy {
		compressors = append(compressors, "snappy")
	}
	if d.conf.Mongo.Compressors.Zlib {
		compressors = append(compressors, "zlib")
	}
	if d.conf.Mongo.Compressors.Zstd {
		compressors = append(compressors, "zstd")
	}
	if len(compressors) > 0 {
		clientOpts.SetCompressors(compressors)
	}

	// TLS
	if d.conf.Mongo.Tls.Use {

		// create the TLS configuration
		tlsConf := NewTlsConfig(d.conf.Mongo.Tls.Crt, d.conf.Mongo.Tls.Key, d.conf.Mongo.Tls.Root, true)
		if err := tlsConf.Create(); err != nil {
			return err
		}

		// apply the TLS configuration to the options
		clientOpts.SetTLSConfig(tlsConf.Config)
	}

	// apply the url
	clientOpts.ApplyURI(d.conf.Mongo.Uri)

	if d.conf.Options.debug {
		log.Printf("INFO: [DATABASE] connecting to Mongo DB server at %s\n", d.conf.Mongo.Uri)
	}

	var err error

	d.Client, err = mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return fmt.Errorf("ERROR: [DATABASE] failed to connect to mongo server at %s REASON: %s", d.conf.Mongo.Uri, err.Error())
	}

	if err = d.Client.Ping(context.TODO(), &readpref.ReadPref{}); err != nil {
		return fmt.Errorf("ERROR: [DATABASE] failed to ping to mongo server at %s REASON: %s", d.conf.Mongo.Uri, err.Error())
	}

	d.Database = d.Client.Database(d.conf.Mongo.Database, &options.DatabaseOptions{})

	log.Printf("INFO: [DATABASE] connected to Mongo DB server at %s\n", d.conf.Mongo.Uri)

	return nil
}

func (d *Database) GetCollection(name string) *mongo.Collection {

	if d.Database != nil {
		return d.Database.Collection(name)
	}
	return nil
}

func (d *Database) Disconnect() {

	if d.Client != nil {
		d.Client.Disconnect(context.TODO())
		d.Database = nil
		d.Client = nil
		log.Println("INFO: [DATABASE] disconnected from MongoDB server")
	}
}
