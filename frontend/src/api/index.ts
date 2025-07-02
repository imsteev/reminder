import axios from "axios";
import { API_BASE_URL } from "../constants";

type WindowWithClerk = Window & {
  Clerk?: {
    session?: {
      getToken(): Promise<string | null>;
    };
  };
};

export const API_URL = `${API_BASE_URL}/api`;

export const initAxios = () => {
  axios.interceptors.request.use(async (config) => {
    config.baseURL = API_URL;
    const token = await getSessionToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });
};

export const getSessionToken = async () => {
  const win = window as WindowWithClerk;
  return (await win.Clerk?.session?.getToken()) ?? null;
};
