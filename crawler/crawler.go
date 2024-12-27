package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type Leiloes struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Price  string `json:"price,optional"`
	Local  string `json:"local,optional"`
	Img    string `json:"img,optional"`
	Tipo   string `json:"tipo,optional"`
	Status string `json:"status,optional"`
}

func GetLeiloesHandler(w http.ResponseWriter, r *http.Request) {
	// Captura o parâmetro "url" da query string
	query := r.URL.Query()
	url := query.Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Parâmetro 'url' é obrigatório"}`))
		return
	}

	// Captura os leilões e envia a resposta
	leiloes, err := getLeiloes(url)
	if err != nil {
		log.Println("Erro ao capturar os leilões:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leiloes)
}

func getLeiloes(startURL string) ([]Leiloes, error) {
	var leiloes []Leiloes
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	// Captura os cards de leilões
	c.OnHTML("div[class=card open]", func(h *colly.HTMLElement) {
		log.Println("Acessei o card:")
		name := h.ChildText("div[class=wrap] a[class=card-title]")
		price := h.ChildText("div[class=card-price]")
		local := h.ChildText("div[class=card-locality]")
		img := h.ChildAttr("a", "data-bg")
		tipo := h.ChildText("div[class=card-instance-title]")
		status := h.ChildText("div[class=card-status]")

		leilao := Leiloes{
			Name:   name,
			URL:    h.Request.URL.String(),
			Tipo:   tipo,
			Img:    img,
			Price:  price,
			Local:  local,
			Status: status,
		}
		leiloes = append(leiloes, leilao)
	})

	// Log de requisições
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visitando:", r.URL.String())
	})

	// Tratamento de erros
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Erro ao visitar:", r.Request.URL, " - ", err)
	})

	// Visita a URL inicial
	err := c.Visit(startURL)
	if err != nil {
		return nil, err
	}

	c.Wait()
	return leiloes, nil
}

func getLeiloesPaginated(baseURL string, maxPages int) ([]Leiloes, error) {
	var leiloes []Leiloes

	for i := 1; i <= maxPages; i++ {
		url := fmt.Sprintf("%s?pagina=%d", baseURL, i)
		pageLeiloes, err := getLeiloes(url)
		if err != nil {
			log.Println("Erro ao capturar página:", i, "-", err)
			continue
		}
		leiloes = append(leiloes, pageLeiloes...)
	}

	return leiloes, nil
}

func getIndividualLeilaoData(links []Leiloes) ([]Leiloes, error) {
	var leiloes []Leiloes

	for _, leilao := range links {
		c := colly.NewCollector(
			colly.AllowURLRevisit(),
		)

		// Captura os dados específicos de cada leilão
		c.OnHTML("div[class=card.open]", func(h *colly.HTMLElement) {
			name := h.ChildText("div[class=wrap] a[class=card-title]")
			price := h.ChildText("div[class=card-price]")
			local := h.ChildText("div[class=card-locality]")
			img := h.ChildAttr("a", "data-bg")
			tipo := h.ChildText("div[class=card-instance-title]")
			status := h.ChildText("div[class=card-status]")

			e := Leiloes{
				Name:   name,
				URL:    h.Response.Request.URL.String(),
				Tipo:   tipo,
				Img:    img,
				Price:  price,
				Local:  local,
				Status: status,
			}

			leiloes = append(leiloes, e)
		})

		c.OnRequest(func(r *colly.Request) {
			log.Println("Visitando:", r.URL.String())
		})

		// Visita cada leilão individualmente
		c.Visit("https://www.megaleiloes.com.br" + leilao.URL)
		c.Wait()
	}

	return leiloes, nil
}
