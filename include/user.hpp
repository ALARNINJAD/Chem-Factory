#ifndef USER_HPP
#define USER_HPP

#include <iostream>
#include <sqlite3.h>
using namespace std;

class user{
    protected:

    public:

        static bool user_login(string username,string password);
        static bool user_register(string username,string password);

        user(){}
        ~user(){}
};

#endif