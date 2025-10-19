package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Product représente un article vendu sur le site
type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Discount    *int
	Stock       int
	Image       string
}

// Méthode pour calculer le prix final si réduction
func (p Product) FinalPrice() float64 {
	if p.Discount != nil {
		return p.Price * (1 - float64(*p.Discount)/100)
	}
	return p.Price
}

// Liste de produits affichés sur la page d’accueil
var products = []Product{
	{1, "Sweat à capuche vert clair", "Sweat confortable en coton bio, coupe oversize.", 49.99, nil, 12, "/assets/image/16A.webp"},
	{2, "Sweat noir graphique", "Sweat noir avec imprimé urbain, édition limitée.", 59.99, func() *int { d := 15; return &d }(), 5, "/assets/image/18A.webp"},
	{3, "Sweat vert bouteille", "Sweat épais à capuche, couleur vert profond.", 54.90, nil, 8, "/assets/image/19A.webp"},
	{4, "Sweat bleu marine", "Sweat à manches longues, coupe droite classique.", 64.00, nil, 10, "/assets/image/21A.webp"},
	{5, "Sweat noir basique", "Sweat minimaliste noir, parfait pour un look casual.", 44.90, nil, 9, "/assets/image/22A.webp"},
	{6, "Pantalon cargo noir", "Cargo streetwear à poches, coupe large.", 69.00, nil, 6, "/assets/image/33B.webp"},
	{7, "Jean large bleu clair", "Jean coupe baggy, teinte bleu clair.", 79.90, nil, 7, "/assets/image/34B.webp"},
}

// Chargement des templates HTML
var templates = template.Must(template.ParseGlob("templates/*.html"))

// Route page d’accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Route d’affichage d’un produit (ex: /product?id=2)
func productHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	for _, p := range products {
		if p.ID == id {
			err := templates.ExecuteTemplate(w, "product.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	// Sert les fichiers statiques (CSS, images…)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Routes principales
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product", productHandler)

	log.Println("Serveur démarré sur : http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
