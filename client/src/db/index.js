import low from 'lowdb';
import LocalStorage from 'lowdb/adapters/LocalStorage';

const adapter = new LocalStorage('db');
const db = low(adapter);

db.defaults({ all: [], allcount: 0 })
    .write()

export { db };