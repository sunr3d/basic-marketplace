package logic

import (
	"testing"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
)

func TestValidateAdInput(t *testing.T) {
	tests := []struct {
		name    string
		input   interfaces.AdInput
		wantErr bool
	}{
		{
			name: "валидное объявление",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Продам велосипед",
					Description: "Почти новый, отличное состояние",
					ImageURL:    "https://example.com/image.jpg",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "минимальная длина заголовка",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "abcd",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "максимальная длина заголовка",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       string(make([]byte, 35)),
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "заголовок на 1 символ меньше минимума",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "abc",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "заголовок на 1 символ больше максимума",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       string(make([]byte, 36)),
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "минимальная длина описания",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "abcd",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "максимальная длина описания",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: string(make([]byte, 100)),
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "описание на 1 символ меньше минимума",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "abc",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "описание на 1 символ больше максимума",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: string(make([]byte, 101)),
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "минимальная цена",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       1.0,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "максимальная цена",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       1000000.0,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "цена меньше минимума",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       0.99,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "цена больше максимума",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       1000000.01,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "валидная ссылка на изображение (jpg)",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.jpg",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "валидная ссылка на изображение (png)",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.png",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
		{
			name: "невалидная ссылка на изображение (bmp)",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "https://example.com/image.bmp",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: true,
		},
		{
			name: "пустая ссылка на изображение (разрешено)",
			input: interfaces.AdInput{
				AdvBase: interfaces.AdvBase{
					Title:       "Заголовок",
					Description: "Описание",
					ImageURL:    "",
					Price:       10000,
				},
				OwnerID: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		err := validateAdInput(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: validateAdInput() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
