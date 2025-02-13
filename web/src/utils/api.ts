import axios from 'axios';

const API_BASE_URL = 'http://localhost:9999/api/notes';

export const fetchNotes = async () => {
	const response = await axios.get(API_BASE_URL);
	return response.data;
};

export const createNote = async (title: string, content: string) => {
	const response = await axios.post(API_BASE_URL, { title, content });
	return response.data;
};

export const deleteNote = async (id: number) => {
	const response = await axios.delete(`${API_BASE_URL}/${id}`);
	return response.data;
};
