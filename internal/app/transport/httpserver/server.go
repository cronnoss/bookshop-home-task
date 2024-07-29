package httpserver

// HTTPServer is a HTTP server for ports.
type HTTPServer struct {
	userService     UserService
	tokenService    TokenService
	bookService     BookService
	categoryService CategoryService
	cartService     CartService
}

// NewHTTPServer creates a new HTTP server for ports.
func NewHTTPServer(userService UserService,
	tokenService TokenService,
	bookService BookService,
	categoryService CategoryService,
	cartService CartService,
) HTTPServer {
	return HTTPServer{
		userService:     userService,
		tokenService:    tokenService,
		bookService:     bookService,
		categoryService: categoryService,
		cartService:     cartService,
	}
}
