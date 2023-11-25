export class User{
    fName: string;
    lName: string;
    username: string;
    password: string;
    email: string;
    address: string;
    type: string;
    isVerified: boolean
    token: string
    refreshToken: string
    constructor(){
        this.fName = ""
        this.lName = ""
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