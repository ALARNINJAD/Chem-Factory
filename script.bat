g++ -c src/sqlite3.c -Iinclude -o sqlite3.o
g++ -o chemfactory.exe main.cpp src/menu.cpp src/user.cpp sqlite3.o -Iinclude -std=c++17