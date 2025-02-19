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


/**
 * search notes
 * @param keyword
 */
export const searchNotes = async (keyword: string) => {
	console.error(keyword, 'keyword')
	const response = await axios.get(`${API_BASE_URL}/search`, {
		params: { keyword },
	});
	return response.data;
};

// 登录
export const login = async (username: string, password: string) => {
	const response = await axios.post(`${API_BASE_URL}/auth/login`, { username, password });
	return response.data;
};
