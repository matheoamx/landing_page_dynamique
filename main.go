package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Structure d’un produit
type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Discount    *int
	Stock       int
	Image       string
}

// Calcul du prix final après remise
func (p Product) FinalPrice() float64 {
	if p.Discount != nil {
		return p.Price * (1 - float64(*p.Discount)/100)
	}
	return p.Price
}

// Données simulées (liste initiale de produits)
var products = []Product{
	{1, "Sweat à capuche vert clair", "Sweat confortable en coton bio, coupe oversize.", 49.99, nil, 12, "/assets/image/16A.webp"},
	{2, "Sweat noir graphique", "Sweat noir avec imprimé urbain, édition limitée.", 59.99, func() *int { d := 15; return &d }(), 5, "/assets/image/18A.webp"},
	{3, "Sweat vert bouteille", "Sweat épais à capuche, couleur vert profond.", 54.90, nil, 8, "/assets/image/19A.webp"},
	{4, "Sweat bleu marine", "Sweat à manches longues, coupe droite classique.", 64.00, nil, 10, "/assets/image/21A.webp"},
	{5, "Sweat noir basique", "Sweat minimaliste noir, parfait pour un look casual.", 44.90, nil, 9, "/assets/image/22A.webp"},
	{6, "Pantalon cargo noir", "Cargo streetwear à poches, coupe large.", 69.00, nil, 6, "/assets/image/33B.webp"},
	{7, "Jean large bleu clair", "Jean coupe baggy, teinte bleu clair.", 79.90, nil, 7, "/assets/image/34B.webp"},
}

// Chargement des templates
var templates = template.Must(template.ParseGlob("templates/*.html"))

// Page d’accueil
func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Page de détail produit
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

// Formulaire d’ajout de produit
func addHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "add.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Traitement du formulaire d’ajout
func addSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	stockStr := r.FormValue("stock")
	discountStr := r.FormValue("discount")

	if name == "" || priceStr == "" || stockStr == "" {
		http.Error(w, "Veuillez remplir tous les champs obligatoires.", http.StatusBadRequest)
		return
	}

	price, _ := strconv.ParseFloat(priceStr, 64)
	stock, _ := strconv.Atoi(stockStr)

	var discount *int
	if discountStr != "" {
		if d, err := strconv.Atoi(discountStr); err == nil {
			discount = &d
		}
	}

	id := len(products) + 1

	newProduct := Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Discount:    discount,
		Stock:       stock,
		Image:       "/assets/image/16A.webp", // image par défaut
	}

	products = append(products, newProduct)

	http.Redirect(w, r, "/product?id="+strconv.Itoa(id), http.StatusSeeOther)
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product", productHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/add/submit", addSubmitHandler)

	log.Println("✅ Serveur démarré sur : http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
