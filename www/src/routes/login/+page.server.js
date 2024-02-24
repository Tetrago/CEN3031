import { BASE_API_PATH } from '$lib/env';
import { redirect } from '@sveltejs/kit';
import { parse } from 'cookie';

/** @type {import('./$types').Actions} */
export const actions = {
    default: async ({ cookies, request, url }) => {
        const form = await request.formData();
        const res = await fetch(`${BASE_API_PATH}/auth/login`, {
            method: 'post',
            mode: 'cors',
            body: JSON.stringify({
                email: form.get("email"),
                password: form.get("password")
            })
        });

        if(res.ok) {
            cookies.set("token", parse(res.headers.getSetCookie()[0]).token, { path: '/' });
        }

        redirect(303, url.searchParams.get('redirectTo') ?? '/');
    }
}