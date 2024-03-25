import { BASE_API_PATH } from "$lib/env";

export async function load({ fetch, locals, params }) {
     let data;
     let id;
     // Check if the token is available
     if (locals.token !== undefined && params.groupName) {
          // Replace spaces with the proper encoding for URLs
          const parts = (params.groupName).split(' ')
          const foo = await fetch(`${BASE_API_PATH}/course/group/${parts[0]}/${parts[1]}`)
          const courseId = await foo.json();
          let now = new Date();
          const before = now.getTime();
          const limit = 20
          // console.log(courseId)
          const res = await fetch(`${BASE_API_PATH}/group/history/${courseId}?limit=${limit}&before=${before}`, {
               method: 'get', 
               credentials: 'include', 
               mode: 'cors', 
             });

          data = await res.json()

          id = courseId
     }
     
     return {
          post: {
               chatHistory: data,
               id: id
          }
     }
}

