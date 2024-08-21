import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'

export type EntityType = {
	id: string
	name: string
	status: string
}
export const sseTesting = createApi({
	reducerPath: 'sseTesting',
	tagTypes: ["Reset"],
	baseQuery: fetchBaseQuery({
		baseUrl: 'http://localhost:8080/',
	}),
	endpoints: (builder) => ({
		getConnect: builder.query<EntityType[], void>({
			query: () => `connect`,
			providesTags: ["Reset"]
		}),
		createEntity: builder.mutation<EntityType[], EntityType[]>({
			query: (data) => {
				return {
					method: 'POST',
					url: 'objects',
					headers: {
						'Content-Type': 'application/json',
					},
					body: data
				}
			},
			async onQueryStarted(_, {queryFulfilled, dispatch}) {
				try {
					const {data} = await queryFulfilled;
					dispatch(sseTesting.util.updateQueryData('getConnect', undefined , (draft) => {
						if (!Array.isArray(draft)) {
							return;
						}
						if (draft) {
							draft.push(...data);
						}
					}));
				} catch (error) {
					console.error('Ошибка при обновлении кеша:', error);
				}
			}
		}),
		resetConnect: builder.query<void, void>({
			query: ()=> 'reset',
			async onQueryStarted(_, {queryFulfilled, dispatch}) {
				try {
					await queryFulfilled;
					dispatch(sseTesting.util.invalidateTags(["Reset"]))
				} catch (error) {
					console.error(error)
				}
			}
		})
	}),
})

export const {useCreateEntityMutation, useGetConnectQuery, useLazyResetConnectQuery} = sseTesting
