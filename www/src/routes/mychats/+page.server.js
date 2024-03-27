

import { BASE_API_PATH } from '$lib/env';

/** @type {import('./$types').PageServerLoad} */
export async function load({ fetch, locals }) {

     let data;

     if(locals.token !== undefined)
     {
          const res = await fetch(`${BASE_API_PATH}/user/groups`, {
               // The API's debug page labels this endpoint as a GET request
               method: 'get',
               credentials: 'include',   
               mode: 'cors',
           });
   
           data = await res.json();
     }

     return {
          groups: data
     }

}