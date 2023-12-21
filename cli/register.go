package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MixinNetwork/bot-api-go-client/v2"
	"github.com/urfave/cli/v2"
)

var botMigrateTIPCmdCli = &cli.Command{
	Name:   "migrate",
	Action: botMigrateTIPCmd,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "keystore,k",
			Usage: "keystore download from https://developers.mixin.one/dashboard",
		},
	},
}

func botMigrateTIPCmd(c *cli.Context) error {
	keystore := c.String("keystore")

	dat, err := os.ReadFile(keystore)
	if err != nil {
		panic(err)
	}
	var app Bot
	err = json.Unmarshal([]byte(dat), &app)
	if err != nil {
		panic(err)
	}

	tipPub, tipPriv, _ := ed25519.GenerateKey(rand.Reader)
	log.Printf("Your tip private seed: %s", hex.EncodeToString(tipPriv.Seed()))

	err = bot.UpdateTipPin(context.Background(), app.Pin, hex.EncodeToString(tipPub), app.PinToken, app.ClientID, app.SessionID, app.PrivateKey)
	if err != nil {
		return fmt.Errorf("bot.UpdateTipPin() => %v", err)
	}

	app.Pin = hex.EncodeToString(tipPriv)
	keystoreRaw, _ := json.Marshal(app)
	log.Printf("your new keystore after migrate: %s", string(keystoreRaw))
	return nil
}

var registerSafeCMDCli = &cli.Command{
	Name:   "register",
	Action: registerSafeCMD,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "keystore,k",
			Usage: "keystore download from https://developers.mixin.one/dashboard",
		},
		&cli.StringFlag{
			Name:  "key,s",
			Usage: "seed for spend private key",
		},
	},
}

func registerSafeCMD(c *cli.Context) error {
	keystore := c.String("keystore")
	seed := c.String("key")

	dat, err := os.ReadFile(keystore)
	if err != nil {
		panic(err)
	}
	var app Bot
	err = json.Unmarshal([]byte(dat), &app)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	method := "GET"
	path := "/safe/me"
	token, err := bot.SignAuthenticationTokenWithoutBody(app.ClientID, app.SessionID, app.PrivateKey, method, path)
	if err != nil {
		return err
	}

	me, err := bot.UserMe(ctx, token)
	if err != nil {
		return err
	}
	if me.HasSafe {
		log.Println("user has registed")
		return nil
	}
	s, err := hex.DecodeString(seed)
	if err != nil {
		panic(err)
	}
	if len(s) != ed25519.SeedSize {
		panic("invalid seed")
	}
	privateKey := ed25519.NewKeyFromSeed(s)
	tipPublic := hex.EncodeToString(privateKey[32:])
	sd := hex.EncodeToString(privateKey.Seed())

	me, err = bot.RegisterSafe(ctx, app.ClientID, tipPublic, sd, app.ClientID, app.SessionID, app.PrivateKey, app.Pin, app.PinToken)
	if err != nil {
		return err
	}
	if me.HasSafe {
		log.Println("user registed")
		return nil
	}

	log.Println("user not registed")
	return nil
}