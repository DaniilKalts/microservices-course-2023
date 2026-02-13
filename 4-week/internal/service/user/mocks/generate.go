package mocks

//go:generate minimock -i github.com/DaniilKalts/microservices-course-2023/4-week/internal/repository.UserRepository -o . -s "_minimock.go"
//go:generate minimock -i github.com/DaniilKalts/microservices-course-2023/4-week/internal/clients/database.TxManager -o . -s "_minimock.go"
