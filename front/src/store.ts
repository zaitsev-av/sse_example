import { configureStore } from '@reduxjs/toolkit'
import {useDispatch, useSelector} from "react-redux";
import {sseTesting} from "./api.ts";

export const store = configureStore({
	reducer: {
		[sseTesting.reducerPath]: sseTesting.reducer
	},
	middleware: (getDefaultMiddleware)=>  getDefaultMiddleware().concat(sseTesting.middleware)
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

export const useAppDispatch = useDispatch.withTypes<AppDispatch>()
export const useAppSelector = useSelector.withTypes<RootState>()
