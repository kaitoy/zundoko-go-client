package runner

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/kaitoy/zundoko-go-client/pkg/client"
	"github.com/kaitoy/zundoko-go-client/pkg/logging"
	"github.com/kaitoy/zundoko-go-client/pkg/model"
	"github.com/kaitoy/zundoko-go-client/pkg/util"
)

// Runner starts a Zundoko Kiyoshi.
type Runner interface {
	// Run starts a Zundoko Kiyoshi.
	Run(intervalMillis time.Duration) error
}

type runner struct {
	cl client.Client
}

// NewRunner creates a Runner instance.
func NewRunner(cl client.Client) Runner {
	return &runner{cl}
}

func (r *runner) Run(intervalMillis time.Duration) error {
	defer logging.GetLogger().Sync()
	logging.GetLogger().Info("Start Zundoko Kiyoshi.")

	for {
		zundokos, err := r.cl.GetZundokos()
		if err != nil {
			return fmt.Errorf("failed to get Zundokos: %w", err)
		}
		if isReadyToKiyoshi(zundokos) {
			break
		}

		word := "Zun"
		if rand.Intn(10) < 5 {
			word = "Doko"
		}
		if err = r.cl.PostZundoko(
			&model.Zundoko{
				Id:     util.NewUUID().String(),
				SaidAt: time.Now(),
				Word:   word,
			},
		); err != nil {
			return fmt.Errorf("failed to create a Zundoko: %w", err)
		}
		fmt.Println(word)

		time.Sleep(intervalMillis * time.Millisecond)
	}

	time.Sleep(intervalMillis * time.Millisecond)

	if err := r.cl.PostKiyoshi(
		&model.Kiyoshi{
			Id:     util.NewUUID().String(),
			SaidAt: time.Now(),
		},
	); err != nil {
		return fmt.Errorf("failed to create a Kiyoshi: %w", err)
	}
	fmt.Println("Ki Yo Shi !")

	return nil
}

func isReadyToKiyoshi(zundokos []model.Zundoko) bool {
	numZundokos := len(zundokos)
	if numZundokos < 5 {
		return false
	}

	sort.Slice(zundokos, func(i, j int) bool {
		return zundokos[i].SaidAt.Before(zundokos[j].SaidAt)
	})

	var words []string
	for _, zd := range zundokos[numZundokos-5:] {
		words = append(words, zd.Word)
	}
	return strings.Join(words, "") == "ZunZunZunZunDoko"
}
