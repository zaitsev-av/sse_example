import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'

export type EntityType = {
	id: string
	name: string
	status: string
}
export const sseTesting = createApi({
	reducerPath: 'sseTesting',
	baseQuery: fetchBaseQuery({
		baseUrl: 'http://localhost:8080/',
	}),
	endpoints: (builder) => ({
		getConnect: builder.query<EntityType[], void>({
			query: () => `connect`,
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
					console.log('data перед dispatch', data);
					dispatch(sseTesting.util.updateQueryData('getConnect', undefined , (draft) => {
						console.log('перед иф');
						console.log('!Array.isArray(draft)', !Array.isArray(draft));
						if (!Array.isArray(draft)) {
							return;
						}
						console.log('data', data);
						if (draft) {
							draft.push(...data);
						}
					}));
				} catch (error) {
					console.error('Ошибка при обновлении кеша:', error);
				}

				// const eventSource = new EventSource(`http://localhost:8080/events`);
				//
				// eventSource.onmessage = (ev: MessageEvent<EntityType>) => {
				// 	console.log(ev)
				// 	if (ev.data) {
				// 		console.log('useSseForDocumentUpdates -> updateQueryData 1');
				// 		dispatch(
				// 			sseTesting.util.updateQueryData('getConnect', undefined, (data) => {
				// 				console.log('useSseForDocumentUpdates -> updateQueryData 2');
				// 				if (!Array.isArray(data)) {
				// 					return;
				// 				}
				// 				if (data) {
				// 					console.log('useSseForDocumentUpdates -> updateQueryData 2');
				// 					const index = data.findIndex((el) => el.id === ev.data.id)
				// 					data[index].status = ev.data.status
				// 				}
				// 			})
				// 		)
				//
				// 		eventSource.close()
				// 	}
				// }

			}
		}),
	}),
})

export const {useCreateEntityMutation, useGetConnectQuery} = sseTesting
