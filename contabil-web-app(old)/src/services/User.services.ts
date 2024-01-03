import { Injectable } from "@angular/core"
import { HttpClient } from '@angular/common/http';
import { Observable } from "rxjs";
import { User } from "src/domains/User";
import {environment} from "../environments/environment";

@Injectable()
export class UserServices {

    url = environment.API_BASE_URL + '/user';

    constructor(private http: HttpClient) {}

    getUsers(): Observable<User[]> {
        return this.http.get<User[]>(`${this.url}/all`);
    }

    login(name: string, password: string): Observable<User> {
        return this.http.post<User>(`${this.url}/login`, {name, password});
    }
}
