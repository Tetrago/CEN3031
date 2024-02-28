import { BASE_API_PATH } from '$lib/env';

/** @type {import('./$types').PageServerLoad} */
export async function load({ fetch, locals }) {
    let data;
    if(locals.token !== undefined)
    {
        const res = await fetch(`${BASE_API_PATH}/user/groups`)
        data = await res.json();
    }

    return {
        groups: data
    }
}