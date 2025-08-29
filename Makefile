all:clean test

test:run_mutex run_context

clean:
	go clean --cache

#TESTS
run_context:
	go test --race my_concurency/internal/mycontext/

run_mutex:
	go test --race my_concurency/internal/mymutexcas/
	go test --race my_concurency/internal/mymutextic/
