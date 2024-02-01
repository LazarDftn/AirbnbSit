export class User{
    ID: string;
    firstName: string;
    lastName: string;
    username: string;
    password: string;
    email: string;
    address: string;
    type: string;
    isVerified: boolean
    token: string
    refreshToken: string
    constructor(){
        this.ID = ""
        this.firstName = ""
        this.lastName = ""
        this.username = ""
        this.password = ""
        this.email = ""
        this.address = ""
        this.type = ""
        this.isVerified = false
        this.token = ""
        this.refreshToken = ""
    }
}