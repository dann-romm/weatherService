NAME = weatherBackend

SRC = cmd/weatherBackend/main.go

all: exec

$(NAME):
	go build -o $(NAME) $(SRC)

run:
	go run $(SRC)

exec: $(NAME)
	./$(NAME)

clean:
	$(RM) $(NAME)

re: clean all

.PHONY: run exec all clean re
