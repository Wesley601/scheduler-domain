clean:
	rm -rf build

agenda-create:
	go build -ldflags="-s -w" -o build/agenda-create lambdas/agenda/create/main.go
agenda-list:
	go build -ldflags="-s -w" -o build/agenda-list lambdas/agenda/list/main.go
agenda-list-slots:
	go build -ldflags="-s -w" -o build/agenda-list-slots lambdas/agenda/list-slots/main.go
agenda-by-id:
	go build -ldflags="-s -w" -o build/agenda-by-id lambdas/agenda/by-id/main.go
booking-book:
	go build -ldflags="-s -w" -o build/book lambdas/booking/book/main.go

agenda: agenda-create agenda-list agenda-list-slots agenda-by-id
booking: booking-book
build: agenda booking
