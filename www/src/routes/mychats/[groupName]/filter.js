// Please Lord forgive me

const pottyWords = ['cunt', 'pussy', 'nigger', 'nigga', 'anal', 'ass', 'asshole',
     'cock', 'bitch', 'bitches', 'blowjob', 'titty', 'clit', 'cum', 'fuck', 'coon', 'coons',
     'cuck', 'fag', 'faggot', 'fucking', 'fuckin', 'motherfucker', 'chink', 'shit', 'tits', 
     'motherfucking', 'chinky', 'coonass', 'dink', 'goombah', 'injun', 'jewboy', 
     'niggeritis', 'raghead', 'dothead', 'sand nigger', 'chingchong', 'towelhead', 'rag head', 
     'dot head', 'towel head', 'wigger', 'wigga', 'whigger', 'white nigger', 'nigger wop', 
     'wog', 'zipper head']


const pattern = pottyWords.map(word => `\\b${word.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&')}\\b`).join('|');
export const restrictedWordsRegex = new RegExp(pattern, 'gi');
