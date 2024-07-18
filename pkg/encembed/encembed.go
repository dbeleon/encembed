package encembed

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"math"
	mrand "math/rand"
	"os"
	"text/template"
	"time"

	"filippo.io/age"
	"github.com/dbeleon/scr"
)

func KeyGen() string {
	//base64 key
	keybytes := make([]byte, 32)
	rand.Read(keybytes)
	return base64.RawStdEncoding.EncodeToString(keybytes)
}

type Config struct {
	PkgName          string
	FuncName         string
	Key              string
	Scramble         bool
	EmbedName        string
	EncryptedVarName string
	DecryptedVarName string
	ExternalKey      string
	ScrLength        int
	ScrFeedbacks     []int
	ScrPolynomial    uint64

	Infile  string
	Outfile string
}

func Embed(cfg Config, byts []byte) error {
	var inf io.Reader
	var err error
	if len(byts) == 0 {
		inf, err = os.Open(cfg.Infile)
		if err != nil {
			return err
		}
	} else {
		inf = bytes.NewReader(byts)
	}

	of, err := os.Create(cfg.EmbedName)
	if err != nil {
		return err
	}

	rcpt, err := age.NewScryptRecipient(cfg.Key)
	if err != nil {
		return err
	}

	w, err := age.Encrypt(of, rcpt)
	if err != nil {
		return err
	}
	defer w.Close()
	io.Copy(w, inf)

	if len(cfg.Outfile) > 0 {
		srcf, err := os.Create(cfg.Outfile)
		if err != nil {
			return err
		}
		defer srcf.Close()

		tmp, err := template.New("encthing").Parse(tpl)
		if err != nil {
			return err
		}

		if cfg.Scramble {
			r := mrand.New(mrand.NewSource(time.Now().Unix()))

			l := int(5 + float64(r.Intn(60)))
			feedbacks := make([]int, int(math.Max(2, float64(r.Intn(int(float32(l)*2.0/3.0)-2)))))
			idxs := make([]int, l)
			for i := 0; i < l; i++ {
				idxs[i] = i
			}

			shuffleInts(r, idxs)
			for i := 0; i < len(feedbacks); i++ {
				feedbacks[i] = idxs[i]
			}
			poly := mrand.Uint64()
			poly = poly & (uint64(0xFFFF_FFFF_FFFF_FFFF) >> (64 - l))
			cfg.ScrLength = l
			cfg.ScrFeedbacks = feedbacks
			cfg.ScrPolynomial = poly

			scram := scr.New(cfg.ScrLength, cfg.ScrFeedbacks, cfg.ScrPolynomial)
			data := []byte(cfg.Key)
			scram.ScrambleAdditive(data)
			cfg.Key = base64.StdEncoding.EncodeToString(data)
		}

		err = tmp.Execute(srcf, cfg)
		if err != nil {
			return err
		}
	}

	if cfg.ExternalKey != "" {
		kf, err := os.Create(cfg.ExternalKey)
		if err != nil {
			return err
		}
		kf.WriteString(cfg.Key)
		kf.Close()
	}
	return nil
}

func shuffleInts(r *mrand.Rand, a []int) {
	for i := range a {
		j := r.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
