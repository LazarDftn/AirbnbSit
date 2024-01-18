export class UserProfile{
    fName: string;
    lName: string;
    username: string;
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
        this.email = ""
        this.address = ""
        this.type = ""
        this.isVerified = false
        this.token = ""
        this.refreshToken = ""
    }
}