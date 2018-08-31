package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "html/template"
  "regexp"
  "errors"
  "sort"
)

type Pagina struct{
  Titulo string
  Cuerpo []byte

  Siguiente string
  Visibilidad_A string
  Anterior string
  Visibilidad_S string
}


var plantillas = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html", "tmpl/front.html"))
var regex_ruta = regexp.MustCompile("^(/|(/(edit|save|view)/([a-zA-Z0-9]+)))$")
var pagina_principal = "Home"
var enlaces []string

func main() {
  http.HandleFunc("/", llamarManejador(manejadorRaiz)) //Nuevo manejador
  http.HandleFunc("/view/", llamarManejador(manejadorMostrar))
  http.HandleFunc("/save/", llamarManejador(manejadorGuardar))
  http.HandleFunc("/edit/", llamarManejador(manejadorEditar))
  http.HandleFunc("/crear/", llamarManejador(manejadorCrear))
  fmt.Println("El servidor se encuentra en ejecución.");
  http.ListenAndServe(":8080", nil)
}

//Método para guardar página
func ( p* Pagina ) guardar() error {
  nombre := p.Titulo + ".txt"

  agregarEnlace(p.Titulo) //Nuevo

  return ioutil.WriteFile( "./view/" + nombre, p.Cuerpo, 0600)
}

//Función para cargar página
func cargarPagina( titulo string ) (*Pagina, error) {
  nombre_archivo := titulo + ".txt"
  fmt.Println("El cliente ha pedido: " + nombre_archivo)
  cuerpo, err := ioutil.ReadFile( "./view/" + nombre_archivo )
  if err != nil {
    return nil, err
  }
  ant, sig := obtenerEnlaces(titulo) //Nuevo
  //Modificado
  return &Pagina{Titulo: titulo, Cuerpo: cuerpo, Anterior: ant, Siguiente: sig}, nil
}

//Función para validar ruta y regresar nombre de la página solicitada
func dameTitulo(w http.ResponseWriter, r *http.Request) (string, error) {
  peticion := regex_ruta.FindStringSubmatch(r.URL.Path)
  if peticion == nil {
    http.NotFound(w, r)
    return "", errors.New("Ruta inválida")
  }
  return peticion[len(peticion) - 1], nil
}

//Función para cargar las plantillas HTML
func cargarPlantilla(w http.ResponseWriter, nombre_plantilla string, pagina *Pagina){
  plantillas.ExecuteTemplate(w, nombre_plantilla + ".html", pagina)
}
func llamarManejador(manejador func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    titulo, err := dameTitulo(w, r)
    fmt.Println(titulo)
    if err != nil {
      http.NotFound(w, r)
      return
    }
    manejador(w, r, titulo)
  }
}

func iniciarEnlaces(){
  paginas, _ := ioutil.ReadDir("./view/")
  enlaces = make([]string, 0, 0)
  for _, actual := range paginas {
    if(actual.Name() != (pagina_principal + ".txt")){
      enlaces = append(enlaces, actual.Name()[ : len(actual.Name()) - 4 ])
    }
  }
}

func agregarEnlace(pagina string){
  for _, actual := range enlaces{
    if (actual == pagina || actual == pagina_principal) {
      return
    }
  }
  enlaces = append(enlaces, pagina)
  sort.Strings(enlaces)
}

func obtenerEnlaces(pagina string) (string, string) {
  var ant, sig string

  for i, actual := range enlaces {
    if(actual == pagina){
      if(i > 0){
        ant = enlaces[i - 1]
      } else{
        ant = ""
      }

      if(i < (len(enlaces)-1) ){
        sig = enlaces[i + 1]
      } else {
        sig = ""
      }
      break
    }
  }
  return ant, sig
}

func ( p* Pagina ) asignarVisibilidad() {
  if(p.Anterior == ""){
    p.Visibilidad_A = "hidden"
  }else {
    p.Visibilidad_A = "visible"
  }
  if(p.Siguiente == ""){
    p.Visibilidad_S = "hidden"
  }else {
    p.Visibilidad_S = "visible"
  }
}

//Manejador para mostrar página principal
func manejadorRaiz(w http.ResponseWriter, r *http.Request, titulo string) {
  p, err := cargarPagina(pagina_principal)
  if err != nil {
    http.Redirect(w, r, "view/" + pagina_principal, http.StatusFound)
    fmt.Println("Error")
    return
  }

  cargarPlantilla(w, "front", p)
}

//Manejador para visualizar wikis
func manejadorMostrar(w http.ResponseWriter, r *http.Request, titulo string){
  p, err := cargarPagina(titulo)
  if err != nil {
    http.Redirect(w, r, "/edit/" + titulo, http.StatusFound)
    fmt.Println("La página solicitada no existía. Llamando al editor...")
    return
  }
  p.asignarVisibilidad() //Nuevo
  cargarPlantilla(w, "view", p)
}

//Manejador para editar wikis
func manejadorEditar(w http.ResponseWriter, r *http.Request, titulo string){
  p, err := cargarPagina(titulo)
  if err != nil{
    p = &Pagina{Titulo: titulo}
  }
  cargarPlantilla(w, "edit", p)
}

//Manejador para guardar wikis
func manejadorGuardar(w http.ResponseWriter, r * http.Request, titulo string) {
  cuerpo := r.FormValue("body")
  p := &Pagina{Titulo: titulo, Cuerpo: []byte(cuerpo)}
  fmt.Println("Guardando " + titulo + ".txt...")
  p.guardar()
  http.Redirect(w, r, "/view/" + titulo, http.StatusFound)
}

//Manejador para crear wikis
func manejadorCrear(w http.RespomseWriter, r* http.Request, titulo string)  {
  p, err := crearPagina()
  if err
}
