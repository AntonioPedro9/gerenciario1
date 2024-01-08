export type IService = {
  id: number;
  name: string;
  description: string;
  duration: number;
  price: number;
  userID: string;
};

export type ICreateServiceRequest = {
  name: string;
  description?: string;
  duration?: number;
  price?: number;
  userID: string;
};

export type IUpdateServiceRequest = {
  id: number;
  name: string;
  description?: string;
  duration?: number;
  price?: number;
  userID: string;
};
