# Simple Chat2Desk API service

Chat2Desk API functions

# Functions

## Initial functions

<details>
<summary>Functions list</summary>

```func (*Ctd).Init(url string, token string)```

<details>
<summary>Function description</summary>

Init initializes the Ctd instance with the provided URL and token.
It sets the URL to ensure it ends with a slash and assigns the token.
The timeout is set to 10 seconds by default.
This method is typically called before making any API requests to ensure
that the Ctd instance is properly configured with the necessary
URL and authentication token.

Parameters:
  - url: The base URL of the Chat2Desk API, which should end with a slash.
  - token: The authentication token for the Chat2Desk API, used
</details>

```func (*Ctd).Get(ctx context.Context, path string, response any) ([]byte, error)```

<details>
<summary>Function description</summary>

Get retrieves data from the specified path using a GET request.
It constructs the full URL by appending the path to the base URL.
The method sends a GET request to the API and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request times out, it retries the request once.
This method is typically used to fetch data from the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - path: The path to the specific API endpoint to retrieve data from.
  - response: A pointer to a struct where the response data will be unmarshaled.

Returns:
  - A byte slice containing the response data from the API.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).Post(ctx context.Context, path string, data any, response any) ([]byte, error)```

<details>
<summary>Function description</summary>

Post sends data to the specified path using a POST request.
It constructs the full URL by appending the path to the base URL.
The method sends a POST request to the API with the provided data and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request times out, it retries the request once.
This method is typically used to send data to the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - path: The path to the specific API endpoint to send data to.
  - data: The data to be sent in the request body, which can be of any type (string, byte slice, or struct).
  - response: A pointer to a struct where the response data will be unmarshaled.

Returns:
  - A byte slice containing the response data from the API.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).Put(ctx context.Context, path string, data any, response any) ([]byte, error)```

<details>
<summary>Function description</summary>

Put sends data to the specified path using a PUT request.
It constructs the full URL by appending the path to the base URL.
The method sends a PUT request to the API with the provided data and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request times out, it retries the request once.
This method is typically used to update data in the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - path: The path to the specific API endpoint to send data to.
  - data: The data to be sent in the request body, which can be of any type (string, byte slice, or struct).
  - response: A pointer to a struct where the response data will be unmarshaled.

Returns:
  - A byte slice containing the response data from the API.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).Delete(ctx context.Context, path string, response any) ([]byte, error)```

<details>
<summary>Function description</summary>

Delete sends a DELETE request to the specified path.
It constructs the full URL by appending the path to the base URL.
The method sends a DELETE request to the API and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request times out, it retries the request once.
This method is typically used to delete data from the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - path: The path to the specific API endpoint to delete data from.
  - response: A pointer to a struct where the response data will be unmarshaled.

Returns:
  - A byte slice containing the response data from the API.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).doRequest(ctx context.Context, method string, url string, payload any, response any) ([]byte, error)```

<details>
<summary>Function description</summary>

doRequest performs an HTTP request with the specified method, URL, and payload.
It handles the request creation, sending, and response reading.
The method supports GET, POST, PUT, and DELETE requests.
It sets the appropriate headers, including the Authorization header if a token is provided.
It also measures the time taken for the request and logs debug information.
If the response body contains an error message indicating an invalid token,
it returns an ErrorInvalidToken error.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - method: The HTTP method to use for the request (e.g., "GET", "POST", "PUT", "DELETE").
  - url: The full URL for the request, including the base URL and any specific path.
  - payload: The data to be sent in the request body, which can be of any type (string, byte slice, or struct).

Returns:
  - A byte slice containing the response data from the API.
  - An error if the request fails, if the response is invalid, or if the response indicates an invalid token.
</details>

```func (*Ctd).LastError() any```

<details>
<summary>Function description</summary>

LastError returns the last error encountered during API requests.
This method is useful for retrieving the last error that occurred,
allowing for error handling or logging in the application.

Returns:
  - The last error encountered during API requests, or nil if no error occurred.
</details>

</details>

## Channels

<details>
<summary>Functions list</summary>

```func (*Ctd).Channels(ctx context.Context, offset int, limit int) (*ChannelsResponse, error)```

<details>
<summary>Function description</summary>

Channels retrieves a list of channels from the Chat2Desk API.
It takes a context, an offset, and a limit as parameters.
The offset is used for pagination, and the limit specifies the maximum
number of channels to return.
It constructs the API endpoint URL with the provided offset and limit,
sends a GET request to the API, and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns the response data.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching channels.
  - limit: The maximum number of channels to return.

Returns:
  - A pointer to a ChannelsResponse struct containing the list of channels and metadata.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetChannels(ctx context.Context, offset int, limit int) ([]ChannelItem, int, error)```

<details>
<summary>Function description</summary>

GetChannels retrieves a list of channels from the Chat2Desk API.
It uses the Channels method to fetch the channels and handles errors.
If the response status is not "success", it logs an error and returns nil.
It returns a slice of ChannelItem, which contains the channels.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching channels.
  - limit: The maximum number of channels to return.

Returns:
  - A slice of ChannelItem containing the channels.
  - The total number of channels available (for pagination).
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Clients

<details>
<summary>Functions list</summary>

```func (*Ctd).APIGetClient(ctx context.Context, id int) (*ClientResponse, error)```

<details>
<summary>Function description</summary>

APIGetClient retrieves a client by its ID from the Chat2Desk API.
It takes a context and the client ID as parameters.
It constructs the API endpoint URL with the provided client ID,
sends a GET request to the API, and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns the response data.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the client to retrieve.

Returns:
  - A pointer to a ClientsResponse struct containing the client data and metadata.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).APIGetClients(ctx context.Context, offset int, limit int, order string, params string) (*ClientsResponse, error)```

<details>
<summary>Function description</summary>

APIGetClients retrieves a list of clients from the Chat2Desk API.
It takes a context, an offset, a limit, an order, and additional parameters as strings.
The offset is used for pagination, the limit specifies the maximum number of clients to return,
the order specifies the sorting order, and params can include additional query parameters.
It constructs the API endpoint URL with the provided parameters,
sends a GET request to the API, and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns the response data.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching clients.
  - limit: The maximum number of clients to return.
  - order: The sorting order for the clients (e.g., "asc", "desc").
  - params: Additional query parameters to include in the request.

Returns:
  - A pointer to a ClientResponse struct containing the list of clients and metadata.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).APICreateClient(ctx context.Context, phone string, transport string, channel_id int, nickname string, assigned_phone string) (*ClientResponse, error)```

<details>
<summary>Function description</summary>

APICreateClient creates a new client in the Chat2Desk API.
It takes a context, phone number, transport type, channel ID, nickname, and assigned phone as parameters.
It constructs the API endpoint URL, prepares the data to be sent in the request,
sends a POST request to the API, and returns the response data as a pointer to ClientsResponse
struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the ClientsResponse struct containing the new client data.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - phone: The phone number of the client to be created.
  - transport: The transport type for the client (e.g., "whatsapp", "telegram").
  - channel_id: The ID of the channel associated with the client.
  - nickname: The nickname of the client (optional).
  - assigned_phone: The assigned phone number for the client (optional).

Returns:
  - A pointer to a ClientsResponse struct containing the new client data.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetClient(ctx context.Context, id int) (*ClientItem, error)```

<details>
<summary>Function description</summary>

GetClient retrieves a client by its ID from the Chat2Desk API.
It takes a context and the client ID as parameters.
It calls the APIGetClient method to fetch the client data.
If the response contains an error or if no client data is found, it returns an error.
If the client is found, it returns a pointer to the ClientItem struct containing the client details.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the client to retrieve.

Returns:
  - A pointer to a ClientItem struct containing the client details.
  - An error if the request fails, if the response is invalid, or if no client data is found.
</details>

```func (*Ctd).GetClientsList(ctx context.Context, offset int, limit int) ([]ClientItem, int, error)```

<details>
<summary>Function description</summary>

GetClients retrieves a list of clients from the Chat2Desk API.
It uses the APIGetClients method to fetch the clients and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a slice of ClientItem, which contains the clients.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching clients.
  - limit: The maximum number of clients to return.

Returns:
  - A slice of ClientItem containing the clients.
  - The total number of clients available (for pagination).
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).CreateClient(ctx context.Context, phone string, transport string, channel_id int, nickname string, assigned_phone string) (*ClientItem, error)```

<details>
<summary>Function description</summary>

CreateClient creates a new client in the Chat2Desk API.
It takes a context, phone number, transport type, channel ID, nickname, and assigned phone as parameters.
It calls the APICreateClient method to create the client and handles errors.
If the response status is not "success", it sets the last error and returns an error.
If the client is created successfully, it returns a pointer to the ClientItem struct containing the client details.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - phone: The phone number of the client to be created.
  - transport: The transport type for the client (e.g., "whatsapp", "telegram").
  - channel_id: The ID of the channel to which the client belongs.
  - nickname: The nickname of the client (optional).
  - assigned_phone: The phone number assigned to the client (optional).

Returns:
  - A pointer to a ClientItem struct containing the client details.
  - An error if the request fails, if the response is invalid, or if the client could not be created.
</details>

</details>

## Companies API Info

<details>
<summary>Functions list</summary>

```func (*Ctd).APICompaniesApiInfo(ctx context.Context) (*CompaniesApiInfoResponse, error)```

<details>
<summary>Function description</summary>

APIGetCompaniesApiInfo retrieves information about the company using the Chat2Desk API.
It constructs the API endpoint URL, sends a GET request to the API,
and returns the response data as a CompaniesApiInfoResponse struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the CompaniesApiInfoResponse struct.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
</details>

```func (*Ctd).CompaniesApiInfo(ctx context.Context) (*CompaniesApiInfoData, error)```

<details>
<summary>Function description</summary>

GetCompaniesApiInfo retrieves information about the company using the Chat2Desk API.
It uses the APICompaniesApiInfo method to fetch the company information and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a CompaniesApiInfoData struct, which contains the company information.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
</details>

</details>

## Custom Client Fields

<details>
<summary>Functions list</summary>

```func (*Ctd).APICustomClientFields(ctx context.Context) (*CustomClientFieldResponse, error)```

<details>
<summary>Function description</summary>

APICustomClientFields retrieves a list of custom client fields from the Chat2Desk API.
It constructs the API endpoint URL, sends a GET request to the API,
and returns the response data as a CustomClientFieldResponse struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the CustomClientFieldResponse struct.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to a CustomClientFieldResponse struct containing the list of custom client fields
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetCustomClientFields(ctx context.Context) ([]CustomClientFieldItem, error)```

<details>
<summary>Function description</summary>

GetCustomClientFields retrieves a list of custom client fields from the Chat2Desk API.
It uses the APICustomClientFields method to fetch the custom client fields and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a slice of CustomClientFieldItem, which contains the custom client fields.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A slice of CustomClientFieldItem containing the custom client fields.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Dialogs

<details>
<summary>Functions list</summary>

```func (*GetDialogsParams).Params() string```

<details>
<summary>Function description</summary>

</details>

```func (*Ctd).APIGetDialogs(ctx context.Context, params *GetDialogsParams) (*DialogsResponse, error)```

<details>
<summary>Function description</summary>

APIGetDialogs retrieves a list of dialogs from the API.
It takes a context and GetDialogsParams, and returns a DialogsResponse or an error.

Parameters:
  - ctx (context.Context): The context for the request.
  - params (*GetDialogsParams): The parameters for filtering and pagination.

Returns:
  - A pointer to a DialogsResponse containing the response data.
  - An error if the request fails.
</details>

```func (*Ctd).APIGetDialog(ctx context.Context, dialog_id int64) (*DialogResponse, error)```

<details>
<summary>Function description</summary>

APIGetDialog retrieves a dialog by its ID from the API.
It takes a context and a dialog ID, and returns a DialogResponse or an error.

Parameters:
  - ctx (context.Context): The context for the request.
  - dialog_id (int64): The ID of the dialog to retrieve.

Returns:
  - A pointer to a DialogResponse containing the response data.
  - An error if the request fails.
</details>

```func (*Ctd).GetDialogs(ctx context.Context, params *GetDialogsParams) ([]DialogItem, int, error)```

<details>
<summary>Function description</summary>

GetDialogs retrieves a list of dialogs.
It takes a context and GetDialogsParams, and returns a slice of DialogItem or an error.

Parameters:
  - ctx (context.Context): The context for the request.
  - params (*GetDialogsParams): The parameters for filtering and pagination.

Returns:
  - A slice of DialogItem containing the dialogs.
  - The total number of dialogs available (for pagination).
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetDialog(ctx context.Context, dialog_id int64) (*DialogItem, error)```

<details>
<summary>Function description</summary>

GetDialog retrieves a dialog by its ID.
It takes a context and a dialog ID, and returns a DialogItem or an error.

Parameters:
  - ctx (context.Context): The context for the request.
  - dialog_id (int64): The ID of the dialog to retrieve.

Returns:
  - A pointer to a DialogItem containing the dialog data.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Messages

<details>
<summary>Functions list</summary>

```func (*Ctd).APISendMessage(ctx context.Context, message *MessagePayload) (*SendMessageResponse, error)```

<details>
<summary>Function description</summary>

APISendMessage sends a message via the API.
It takes a context and a MessagePayload, and returns a MessageResponse or an error.

Parameters:
  - ctx (context.Context): The context for the request.
  - message (*MessagePayload): The message payload to send.

Returns:
  - A pointer to a MessageResponse containing the response data.
  - An error if the request fails.
</details>

```func (*Ctd).SendMessage(ctx context.Context, message *MessagePayload) (*SendMessageItem, error)```

<details>
<summary>Function description</summary>

SendMessage sends a message to the API.
It takes a context and a MessagePayload, and returns a MessageItem or an error.

Parameters:
  - ctx (context.Context): The context for the request.
  - message (MessagePayload): The message payload to send.

Returns:
  - A pointer to a MessageItem containing the response data.
  - An error if the request fails.
</details>

</details>

## Operator Groups

<details>
<summary>Functions list</summary>

```func (*Ctd).APIOperatorGroups(ctx context.Context) (*OperatorGroupsResponse, error)```

<details>
<summary>Function description</summary>

APIOperatorGroups retrieves a list of operator groups from the Chat2Desk API.
It constructs the API endpoint URL, sends a GET request to the API,
and returns the response data as an OperatorGroupsResponse struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the OperatorGroupsResponse struct.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to an OperatorGroupsResponse struct containing the list of operator groups
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).OperatorGroups(ctx context.Context) ([]OperatorGroup, error)```

<details>
<summary>Function description</summary>

OperatorGroups retrieves a list of operator groups from the Chat2Desk API.
It uses the APIOperatorGroups method to fetch the operator groups and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a slice of OperatorGroup, which contains the operator groups.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A slice of OperatorGroup containing the list of operator groups.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Operators

<details>
<summary>Functions list</summary>

```func (*Ctd).APIOperators(ctx context.Context, offset int, limit int) (*OperatorsResponse, error)```

<details>
<summary>Function description</summary>

APIOperators retrieves a list of operators from the Chat2Desk API.
It constructs the API endpoint URL with the provided offset and limit,
sends a GET request to the API, and returns the response data as a OperatorsResponse struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the OperatorsResponse struct.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching operators.
  - limit: The maximum number of operators to return.

Returns:
  - A pointer to a CtdOperatorsResponse struct containing the list of operators and metadata.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).Operators(ctx context.Context, offset int, limit int) ([]Operator, int, error)```

<details>
<summary>Function description</summary>

Operators retrieves a list of operators from the Chat2Desk API.
It uses the APIOperators method to fetch the operators and handles errors.
If the response status is not "success", it returns nil.
It returns a slice of Operator, which contains the operators.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching operators.
  - limit: The maximum number of operators to return.

Returns:
  - A a slice of Operator containing the list of operators.
  - The total number of operators available (for pagination).
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).AllOperators(ctx context.Context) ([]Operator, error)```

<details>
<summary>Function description</summary>

AllOperators retrieves all operators from the Chat2Desk API by handling pagination.
It repeatedly calls the Operators method with increasing offsets until all operators are fetched.
It returns a slice of Operator, which contains all the operators.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A slice of Operator containing all the operators.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Statistics

<details>
<summary>Functions list</summary>

```func (*StatisticsRating).GetScoreValue() int64```

<details>
<summary>Function description</summary>

GetScoreValue converts the ScoreValue from json.Number to int64.
If the conversion fails, it returns -1 to indicate an invalid score.

Returns:
  - An int64 representing the score value, or -1 if the conversion fails.
</details>

```func (*StatisticsRating).GetRangeValue(limit1 int64, limit2 int64) uint8```

<details>
<summary>Function description</summary>

GetRangeValue categorizes the score value into three ranges based on the provided limits.
It uses the GetScoreValue method to retrieve the score value.
If the score value is less than or equal to limit1, it returns 1.
If the score value is less than or equal to limit2, it returns 2.
If the score value is greater than limit2, it returns 3.
If the score value is invalid (i.e., -1), it returns 0.

Parameters:
  - limit1: The first limit for categorization.
  - limit2: The second limit for categorization.

Returns:
  - A uint8 representing the category of the score value (0, 1, 2, or 3).
</details>

```func (*Ctd).APIStatisticsRating(ctx context.Context, date time.Time, offset int, limit int) (*StatisticsRatingResponse, error)```

<details>
<summary>Function description</summary>

APIStatisticsRating retrieves a list of statistic ratings from the Chat2Desk API.
It constructs the API endpoint URL with the provided date, offset, and limit,
sends a GET request to the API, and returns the response data as a StatisticsRatingResponse struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the StatisticsRatingResponse struct.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - date: The date for which to retrieve statistics. If zero, the current date is used.
  - offset: The offset for pagination, indicating where to start fetching ratings.
  - limit: The maximum number of ratings to return.

Returns:
  - A pointer to a StatisticsRatingResponse struct containing the list of statistic ratings and metadata.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).StatisticsRating(ctx context.Context, date time.Time, offset int, limit int) ([]StatisticsRating, int, error)```

<details>
<summary>Function description</summary>

StatisticsRating retrieves a list of statistic ratings from the Chat2Desk API.
It uses the APIStatisticsRating method to fetch the ratings and handles errors.
If the response status is not "success", it returns nil.
It returns a slice of StatisticsRating, which contains the ratings.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - date: The date for which to retrieve statistics. If zero, the current date is used.
  - offset: The offset for pagination, indicating where to start fetching ratings.
  - limit: The maximum number of ratings to return.

Returns:
  - A slice of StatisticsRating containing the list of statistic ratings.
  - The total number of ratings available (for pagination).
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).AllStatisticsRating(ctx context.Context, date time.Time) ([]StatisticsRating, error)```

<details>
<summary>Function description</summary>

AllStatisticsRating retrieves all statistic ratings from the Chat2Desk API by handling pagination.
It repeatedly calls the StatisticsRating method with increasing offsets until all ratings are fetched.
It returns a slice of StatisticsRating, which contains all the ratings.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - date: The date for which to retrieve statistics. If zero, the current date is used.

Returns:
  - A slice of StatisticsRating containing all the statistic ratings.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Tags

<details>
<summary>Functions list</summary>

```func (*Ctd).APIGetTags(ctx context.Context, offset int, limit int) (*TagsResponse, error)```

<details>
<summary>Function description</summary>

GetTags retrieves a list of tags from the Chat2Desk API.
It uses the APIGetTags method to fetch the tags and handles errors.
It returns a pointer to a slice of TagItem, which contains the tags.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination.
  - limit: The maximum number of tags to retrieve.

Returns:
  - A pointer to a TagsResponse struct containing the list of tags
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).APIGetTag(ctx context.Context, id int) (*TagResponse, error)```

<details>
<summary>Function description</summary>

APIGetTag retrieves a specific tag by its ID from the Chat2Desk API.
It uses the doRequest method to send a GET request to the API.
If the request fails, it logs the error and returns nil.
It returns a pointer to a TagResponse struct containing the tag data.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the tag to retrieve.

Returns:
  - A pointer to a TagResponse struct containing the tag data
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetTags(ctx context.Context, offset int, limit int) ([]TagItem, int, error)```

<details>
<summary>Function description</summary>

GetTags retrieves a list of tags from the Chat2Desk API.
It uses the APIGetTags method to fetch the tags and handles errors.
If the response status is not "success", it returns nil.
It returns a slice of TagItem, which contains the tags.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination.
  - limit: The maximum number of tags to retrieve.

Returns:
  - A slice of TagItem, which contains the tags
  - The total number of tags available.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetTag(ctx context.Context, id int) (*TagItem, error)```

<details>
<summary>Function description</summary>

GetTag retrieves a specific tag by its ID from the Chat2Desk API.
It uses the APIGetTag method to fetch the tag and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a TagItem, which contains the tag data.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the tag to retrieve.

Returns:
  - A pointer to a TagItem, which contains the tag data
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetAllTags(ctx context.Context) ([]TagItem, error)```

<details>
<summary>Function description</summary>

GetAllTags retrieves all tags from the Chat2Desk API.
It uses the GetTags method to fetch tags in a loop until all tags are retrieved.
It returns a slice of TagItem, which contains all the tags.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A slice of TagItem, which contains all the tags.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## WebHooks

<details>
<summary>Functions list</summary>

```func (*CreateWebhookResponse).Error() string```

<details>
<summary>Function description</summary>

Error compiles the error messages from the CreateWebhookResponse into a single string.
It checks the Errors field for any errors related to the URL, order, or events,
and concatenates them into a single string separated by semicolons.
If there are no errors, it returns an empty string.
</details>

```func (*CreateWebhookResponse).Postprocess() error```

<details>
<summary>Function description</summary>

Postprocess processes the response from the CreateWebhook API endpoint.
It checks the status of the response and returns an error if the status is not "success".
</details>

```func (*WebhookPayload).Prepare()```

<details>
<summary>Function description</summary>

Prepare normalizes the status field of the WebhookPayload.
It ensures that the status is set to either "enable" or "disable".
If the status is not one of these values, it defaults to "enable".
This method is typically used to ensure that the status field is in a valid format
before sending the payload to the API.
It is called before creating or updating a webhook to ensure consistency.
It is used to prepare the payload for API requests.
</details>

```func (*Ctd).Webhooks(ctx context.Context) (*WebhooksResponse, error)```

<details>
<summary>Function description</summary>

Webhooks retrieves a list of webhooks from the Chat2Desk API.
It takes a context as a parameter and constructs the API endpoint URL.
It sends a GET request to the API and returns the response data as a byte slice.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it unmarshals the response data into a WebhooksResponse
struct and returns it.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to a WebhooksResponse struct containing the list of webhooks and status.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).PostWebhooks(ctx context.Context, payload *WebhookPayload) (*CreateWebhookResponse, error)```

<details>
<summary>Function description</summary>

PostWebhook creates a new webhook in the Chat2Desk API.
It takes a context and a WebhookPayload as parameters.
It constructs the API endpoint URL and sends a POST request with the payload.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it unmarshals the response data into a WebhookResponse
struct and returns it.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - payload: The WebhookPayload containing the details of the webhook to be created.

Returns:
  - A pointer to a WebhookResponse struct containing the created webhook and status.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).PutWebhooks(ctx context.Context, id int, payload *WebhookPayload) (*CreateWebhookResponse, error)```

<details>
<summary>Function description</summary>

PutWebhooks updates an existing webhook in the Chat2Desk API.
It takes a context, the webhook ID, and a WebhookPayload as parameters.
It constructs the API endpoint URL with the webhook ID and sends a PUT request with the payload.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it unmarshals the response data into a CreateWebhookResponse
struct and returns it.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the webhook to be updated.
  - payload: The WebhookPayload containing the updated details of the webhook.

Returns:
  - A pointer to a CreateWebhookResponse struct containing the updated webhook and status.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).DeleteWebhooks(ctx context.Context, id int) (*DeleteWebhookResponse, error)```

<details>
<summary>Function description</summary>

DeleteWebhooks deletes a webhook in the Chat2Desk API.
It takes a context and the webhook ID as parameters.
It constructs the API endpoint URL with the webhook ID and sends a DELETE request.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it unmarshals the response data into a DeleteWebhookResponse
struct and returns it.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the webhook to be deleted.

Returns:
  - A pointer to a DeleteWebhookResponse struct containing the status of the delete operation.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).GetWebhooks(ctx context.Context) ([]WebhookItem, error)```

<details>
<summary>Function description</summary>

GetWebhooks retrieves a list of webhooks from the Chat2Desk API.
It takes a context as a parameter and calls the Webhooks method.
If the response status is not "success", it logs an error and returns nil.
It returns a slice of WebhookItem, which contains the webhooks.
If an error occurs during the request, it returns nil and the error.
If the request is successful, it returns a pointer to a slice of WebhookItem.
This method is typically used to fetch webhooks from the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A slice of WebhookItem containing the webhooks.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).CreateWebhook(ctx context.Context, payload *WebhookPayload) (*WebhookItem, error)```

<details>
<summary>Function description</summary>

CreateWebhook creates a new webhook in the Chat2Desk API.
It takes a context and a WebhookPayload as parameters.
It calls the PostWebhook method to send the request.
If the response status is not "success", it logs an error and returns nil.
If the URL is already used, it returns an error indicating that the URL is already used.
If the request is successful, it returns a pointer to the created WebhookItem.
This method is typically used to create new webhooks in the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - payload: The WebhookPayload containing the details of the webhook to be created.

Returns:
  - A pointer to a WebhookItem containing the created webhook.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).UpdateWebhook(ctx context.Context, id int, payload *WebhookPayload) (*WebhookItem, error)```

<details>
<summary>Function description</summary>

UpdateWebhook updates an existing webhook in the Chat2Desk API.
It takes a context, the webhook ID, and a WebhookPayload as parameters.
It calls the PutWebhooks method to send the request.
If the response status is not "success", it logs an error and returns nil.
If the URL is already used, it returns an error indicating that the URL is already used.
If the request is successful, it returns a pointer to the updated WebhookItem.
This method is typically used to update existing webhooks in the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the webhook to be updated.
  - payload: The WebhookPayload containing the updated details of the webhook.

Returns:
  - A pointer to a WebhookItem containing the updated webhook.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).DeleteWebhook(ctx context.Context, id int) error```

<details>
<summary>Function description</summary>

DeleteWebhook deletes a webhook in the Chat2Desk API.
It takes a context and the webhook ID as parameters.
It calls the DeleteWebhooks method to send the request.
If the response status is not "success", it logs an error and returns an error.
If the request is successful, it returns nil.
This method is typically used to delete webhooks in the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - id: The ID of the webhook to be deleted.

Returns:
  - An error if the request fails or if the response is invalid.
</details>

</details>



# Used libraries
* https://github.com/ra-company/env - Simple environment library (GPL-3.0 license)
* https://github.com/ra-company/logging - Simple logging library (GPL-3.0 license)
* https://github.com/stretchr/testify - Module for tests (MIT License)
* https://github.com/brianvoe/gofakeit/ - Random data generator written in go (MIT License)

# Staying up to date
To update library to the latest version, use go get -u github.com/ra-company/ctd.

# Supported go versions
We currently support the most recent major Go versions from 1.24.5 onward.

# License
This project is licensed under the terms of the GPL-3.0 license.