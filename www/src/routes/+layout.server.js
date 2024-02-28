import { jwtDecode } from 'jwt-decode';

/** @type {import('./$types').LayoutServerLoad} */
export async function load({ locals }) {
    let ident = locals.token !== undefined ? jwtDecode(locals.token).ident : "";
    return { ident };
}