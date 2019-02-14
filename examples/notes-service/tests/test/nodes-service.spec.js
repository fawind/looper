const expect = require('expect');
const axios = require('axios');

const ENDPOINT = 'http://' + (process.env.ENDPOINT || '0.0.0.0:8080');
console.log('Testing against endpoint:', ENDPOINT)

describe('Notes-Service', () => {

  it('shoud initially fetch empty notes', (done) => {
    const notes = getNotes()
      .then(response => {
        expect(response.status).toBe(200);
        expect(response.data).toEqual([]);
        done();
      })
      .catch(err => done(err));
  });

  it('should add a new note', (done) => {
    const testContent = 'Test content 1';
    const notes = addNote(testContent)
      .then(response => {
        expect(response.status).toBe(200);
        expect(typeof response.data.id).toBe('number');
        expect(response.data.content).toBe(testContent);
        done();
      })
      .catch(err => done(err));
  })

  it('should list added notes', (done) => {
    const notes = getNotes()
      .then(response => {
        expect(response.status).toBe(200);
        expect(response.data.length).toBeGreaterThan(0);
        done();
      })
      .catch(err => done(err));
  })
});

const getNotes = () => axios.get(ENDPOINT + '/notes');
const addNote = (content) => axios.post(ENDPOINT + '/notes', { content });
