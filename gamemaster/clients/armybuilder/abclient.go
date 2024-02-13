package armybuilder

import (
	"context"
	"log"

	"github.com/corstijank/mekstrike/domain/unit"
	"github.com/corstijank/mekstrike/gamemaster/clients/armybuilder/abprotos"
	dapr "github.com/dapr/go-sdk/client"
	"google.golang.org/protobuf/proto"
)

type AbClient struct {
	client dapr.Client
}

func New(client dapr.Client) (AbClient, error) {
	return AbClient{
		client: client,
	}, nil
}

func (c *AbClient) CreateArmy(ctx context.Context, lights int, mediums int, heavies int, assaults int) ([]*unit.Stats, error) {
	freq := &abprotos.ArmyRequest{
		Lights:   int32(lights),
		Mediums:  int32(mediums),
		Heavies:  int32(heavies),
		Assaults: int32(assaults),
	}
	out, err := proto.Marshal(freq)
	if err != nil {
		log.Println(err)
		return make([]*unit.Stats, 0), err
	}

	content := &dapr.DataContent{
		ContentType: "application/protobuf;proto=net.mekstrike.armybuilder.ArmyRequest",
		Data:        out,
	}
	resp, err := c.client.InvokeMethodWithContent(
		ctx,
		"armybuilder",
		"createArmy",
		"GET",
		content)
	if err != nil {
		log.Println(err)
		return make([]*unit.Stats, 0), err

	}

	army := abprotos.ArmyResponse{}
	err = proto.Unmarshal(resp, &army)
	if err != nil {
		log.Println(err)
	}
	return army.Units, nil
}
