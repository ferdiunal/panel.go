// Package openapi, statik endpoint'ler için OpenAPI spesifikasyonu oluşturur.
//
// Bu paket, Panel.go'nun statik endpoint'leri (auth, init, navigation) için
// OpenAPI path ve operation tanımlamalarını oluşturur.
package openapi

// StaticSpecGenerator, statik endpoint'ler için OpenAPI spesifikasyonu oluşturur.
//
// ## Statik Endpoint'ler
//   - Authentication: sign-in, sign-up, sign-out, forgot-password, session
//   - System: init, navigation
//
// ## Kullanım Örneği
//
//	generator := NewStaticSpecGenerator()
//	paths := generator.GenerateStaticPaths()
type StaticSpecGenerator struct{}

// NewStaticSpecGenerator, yeni bir StaticSpecGenerator oluşturur.
//
// ## Dönüş Değeri
//   - *StaticSpecGenerator: Yapılandırılmış generator
func NewStaticSpecGenerator() *StaticSpecGenerator {
	return &StaticSpecGenerator{}
}

// GenerateStaticPaths, statik endpoint'ler için OpenAPI path'leri oluşturur.
//
// ## Dönüş Değeri
//   - map[string]PathItem: Path -> PathItem mapping'i
//
// ## Oluşturulan Path'ler
//   - POST /api/auth/sign-in/email
//   - POST /api/auth/sign-up/email
//   - POST /api/auth/sign-out
//   - POST /api/auth/forgot-password
//   - GET /api/auth/session
//   - GET /api/init
//   - GET /api/navigation
func (g *StaticSpecGenerator) GenerateStaticPaths() map[string]PathItem {
	paths := make(map[string]PathItem)

	// Authentication endpoints
	paths["/api/auth/sign-in/email"] = g.generateSignInPath()
	paths["/api/auth/sign-up/email"] = g.generateSignUpPath()
	paths["/api/auth/sign-out"] = g.generateSignOutPath()
	paths["/api/auth/forgot-password"] = g.generateForgotPasswordPath()
	paths["/api/auth/session"] = g.generateSessionPath()

	// System endpoints
	paths["/api/init"] = g.generateInitPath()
	paths["/api/navigation"] = g.generateNavigationPath()

	return paths
}

// generateSignInPath, sign-in endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - POST /api/auth/sign-in/email
//
// ## Request Body
//   - email: string (required)
//   - password: string (required)
//
// ## Responses
//   - 200: Successful login
//   - 401: Invalid credentials
//   - 429: Too many requests (rate limit)
func (g *StaticSpecGenerator) generateSignInPath() PathItem {
	return PathItem{
		Post: &Operation{
			Summary:     "Email ile giriş yap",
			Description: "Kullanıcı email ve şifre ile giriş yapar. Başarılı girişte session cookie oluşturulur.",
			OperationID: "signInEmail",
			Tags:        []string{"auth"},
			RequestBody: &RequestBody{
				Description: "Giriş bilgileri",
				Required:    true,
				Content: map[string]MediaType{
					"application/json": {
						Schema: &Schema{
							Type: "object",
							Properties: map[string]Schema{
								"email": {
									Type:        "string",
									Format:      "email",
									Description: "Kullanıcı email adresi",
									Example:     "user@example.com",
								},
								"password": {
									Type:        "string",
									Format:      "password",
									Description: "Kullanıcı şifresi",
									Example:     "********",
								},
							},
							Required: []string{"email", "password"},
						},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Başarılı giriş",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"user": {
										Type: "object",
										Properties: map[string]Schema{
											"id":    {Type: "integer", Example: 1},
											"name":  {Type: "string", Example: "John Doe"},
											"email": {Type: "string", Format: "email", Example: "user@example.com"},
										},
									},
								},
							},
						},
					},
				},
				"401": {
					Description: "Geçersiz kimlik bilgileri",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/ErrorResponse"},
						},
					},
				},
				"429": {
					Description: "Çok fazla istek (rate limit)",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/ErrorResponse"},
						},
					},
				},
			},
			Security: []SecurityRequirement{}, // No authentication required
		},
	}
}

// generateSignUpPath, sign-up endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - POST /api/auth/sign-up/email
//
// ## Request Body
//   - name: string (required)
//   - email: string (required)
//   - password: string (required)
//   - password_confirmation: string (required)
//
// ## Responses
//   - 201: User created
//   - 400: Validation error
//   - 429: Too many requests
func (g *StaticSpecGenerator) generateSignUpPath() PathItem {
	return PathItem{
		Post: &Operation{
			Summary:     "Email ile kayıt ol",
			Description: "Yeni kullanıcı kaydı oluşturur. Kayıt özelliği aktif olmalıdır.",
			OperationID: "signUpEmail",
			Tags:        []string{"auth"},
			RequestBody: &RequestBody{
				Description: "Kayıt bilgileri",
				Required:    true,
				Content: map[string]MediaType{
					"application/json": {
						Schema: &Schema{
							Type: "object",
							Properties: map[string]Schema{
								"name": {
									Type:        "string",
									Description: "Kullanıcı adı",
									Example:     "John Doe",
								},
								"email": {
									Type:        "string",
									Format:      "email",
									Description: "Kullanıcı email adresi",
									Example:     "user@example.com",
								},
								"password": {
									Type:        "string",
									Format:      "password",
									Description: "Kullanıcı şifresi (minimum 8 karakter)",
									MinLength:   ptr(8),
									Example:     "********",
								},
								"password_confirmation": {
									Type:        "string",
									Format:      "password",
									Description: "Şifre tekrarı",
									Example:     "********",
								},
							},
							Required: []string{"name", "email", "password", "password_confirmation"},
						},
					},
				},
			},
			Responses: map[string]Response{
				"201": {
					Description: "Kullanıcı oluşturuldu",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"user": {
										Type: "object",
										Properties: map[string]Schema{
											"id":    {Type: "integer", Example: 1},
											"name":  {Type: "string", Example: "John Doe"},
											"email": {Type: "string", Format: "email", Example: "user@example.com"},
										},
									},
								},
							},
						},
					},
				},
				"400": {
					Description: "Validasyon hatası",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/ErrorResponse"},
						},
					},
				},
				"429": {
					Description: "Çok fazla istek",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/ErrorResponse"},
						},
					},
				},
			},
			Security: []SecurityRequirement{}, // No authentication required
		},
	}
}

// generateSignOutPath, sign-out endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - POST /api/auth/sign-out
//
// ## Responses
//   - 200: Successfully signed out
func (g *StaticSpecGenerator) generateSignOutPath() PathItem {
	return PathItem{
		Post: &Operation{
			Summary:     "Çıkış yap",
			Description: "Kullanıcı oturumunu sonlandırır ve session cookie'yi siler.",
			OperationID: "signOut",
			Tags:        []string{"auth"},
			Responses: map[string]Response{
				"200": {
					Description: "Başarıyla çıkış yapıldı",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/SuccessResponse"},
						},
					},
				},
			},
		},
	}
}

// generateForgotPasswordPath, forgot-password endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - POST /api/auth/forgot-password
//
// ## Request Body
//   - email: string (required)
//
// ## Responses
//   - 200: Password reset email sent
//   - 429: Too many requests
func (g *StaticSpecGenerator) generateForgotPasswordPath() PathItem {
	return PathItem{
		Post: &Operation{
			Summary:     "Şifremi unuttum",
			Description: "Şifre sıfırlama bağlantısı içeren email gönderir.",
			OperationID: "forgotPassword",
			Tags:        []string{"auth"},
			RequestBody: &RequestBody{
				Description: "Email adresi",
				Required:    true,
				Content: map[string]MediaType{
					"application/json": {
						Schema: &Schema{
							Type: "object",
							Properties: map[string]Schema{
								"email": {
									Type:        "string",
									Format:      "email",
									Description: "Kullanıcı email adresi",
									Example:     "user@example.com",
								},
							},
							Required: []string{"email"},
						},
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Şifre sıfırlama email'i gönderildi",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/SuccessResponse"},
						},
					},
				},
				"429": {
					Description: "Çok fazla istek",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/ErrorResponse"},
						},
					},
				},
			},
			Security: []SecurityRequirement{}, // No authentication required
		},
	}
}

// generateSessionPath, session endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - GET /api/auth/session
//
// ## Responses
//   - 200: Session information
//   - 401: Not authenticated
func (g *StaticSpecGenerator) generateSessionPath() PathItem {
	return PathItem{
		Get: &Operation{
			Summary:     "Oturum bilgisi al",
			Description: "Mevcut kullanıcının oturum bilgilerini döndürür.",
			OperationID: "getSession",
			Tags:        []string{"auth"},
			Responses: map[string]Response{
				"200": {
					Description: "Oturum bilgisi",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"user": {
										Type: "object",
										Properties: map[string]Schema{
											"id":    {Type: "integer", Example: 1},
											"name":  {Type: "string", Example: "John Doe"},
											"email": {Type: "string", Format: "email", Example: "user@example.com"},
										},
									},
								},
							},
						},
					},
				},
				"401": {
					Description: "Kimlik doğrulanmadı",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{Ref: "#/components/schemas/ErrorResponse"},
						},
					},
				},
			},
			Security: []SecurityRequirement{}, // No authentication required (returns 401 if not authenticated)
		},
	}
}

// generateInitPath, init endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - GET /api/init
//
// ## Responses
//   - 200: Application initialization data
func (g *StaticSpecGenerator) generateInitPath() PathItem {
	return PathItem{
		Get: &Operation{
			Summary:     "Uygulama başlatma bilgileri",
			Description: "Uygulamanın başlatılması için gerekli bilgileri döndürür (özellikler, OAuth ayarları, versiyon).",
			OperationID: "getInit",
			Tags:        []string{"system"},
			Responses: map[string]Response{
				"200": {
					Description: "Başlatma bilgileri",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"features": {
										Type: "object",
										Properties: map[string]Schema{
											"register":        {Type: "boolean", Description: "Kayıt özelliği aktif mi?", Example: true},
											"forgot_password": {Type: "boolean", Description: "Şifremi unuttum özelliği aktif mi?", Example: false},
										},
									},
									"oauth": {
										Type: "object",
										Properties: map[string]Schema{
											"google": {Type: "boolean", Description: "Google OAuth aktif mi?", Example: false},
										},
									},
									"version": {
										Type:        "string",
										Description: "API versiyonu",
										Example:     "1.0.0",
									},
									"settings": {
										Type:        "object",
										Description: "Dinamik ayarlar",
									},
								},
							},
						},
					},
				},
			},
			Security: []SecurityRequirement{}, // No authentication required
		},
	}
}

// generateNavigationPath, navigation endpoint'i için PathItem oluşturur.
//
// ## Endpoint
//   - GET /api/navigation
//
// ## Responses
//   - 200: Navigation menu items
func (g *StaticSpecGenerator) generateNavigationPath() PathItem {
	return PathItem{
		Get: &Operation{
			Summary:     "Navigasyon menüsü",
			Description: "Yan menü (sidebar) için navigasyon öğelerini döndürür.",
			OperationID: "getNavigation",
			Tags:        []string{"system"},
			Responses: map[string]Response{
				"200": {
					Description: "Navigasyon öğeleri",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"data": {
										Type: "array",
										Items: &Schema{
											Type: "object",
											Properties: map[string]Schema{
												"slug":  {Type: "string", Description: "Resource/page slug", Example: "users"},
												"title": {Type: "string", Description: "Görüntülenecek başlık", Example: "Kullanıcılar"},
												"icon":  {Type: "string", Description: "İkon adı", Example: "users"},
												"group": {Type: "string", Description: "Menü grubu", Example: "Yönetim"},
												"type":  {Type: "string", Description: "Tip (resource veya page)", Example: "resource"},
												"order": {Type: "integer", Description: "Sıralama", Example: 1},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
