#include "../include/menu.hpp"
#include "../include/user.hpp"

void menu::auth_menu(){
    cout << "1) login" << endl
         << "2) register" << endl
         << "choose an option (1-2) -> ";
    char choice;
    string username,password;
    cin >> choice;
    switch (choice){
        case '1':
            cout << "username: ";
            cin >> username;
            cout << "password: ";
            cin >> password;
            if user_login(username,password){
                cout << "logged in!" << endl;
                
            } else {
                cout << "failed!";
            }
            break;
        
        case '2':
            cout << "username: ";
            cin >> username;
            cout << "password: ";
            cin >> password;
            if user_register(username,password){
                cout << "registered!" << endl;
                
            } else {
                cout << "failed!" << endl;
            }
            break;

        default:
            cout << "wrong choice!" << endl;
            break;
    }
}

void menu::game_menu(){
    cout << "1) mix" << endl
         << "2) shop" << endl
         << "3) customers" << endl;
}