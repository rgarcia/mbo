package main

func main() {
	Parse(
		new(Login),
		new(Logout),
		new(Ls),
		new(Register),
		new(Schedule),
	)
}
