# Simple Chat2Desk API service

Chat2Desk API functions

# Functions

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

```func (*Ctd).GetChannels(ctx context.Context, offset int, limit int) (*[]ChannelItem, error)```

<details>
<summary>Function description</summary>

GetChannels retrieves a list of channels from the Chat2Desk API.
It uses the Channels method to fetch the channels and handles errors.
If the response status is not "success", it logs an error and returns nil.
It returns a pointer to a slice of ChannelItem, which contains the channels.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching channels.
  - limit: The maximum number of channels to return.

Returns:
  - A pointer to a slice of ChannelItem containing the channels.
  - An error if the request fails or if the response is invalid.

This function is a wrapper around the Channels method to provide a more user-friendly interface.
It simplifies the process of fetching channels by handling the response and error checking.
It is useful for applications that need to retrieve channels from the Chat2Desk API
in a straightforward manner without dealing with the raw response data.
It is designed to be used in contexts where channels need to be displayed or processed further.
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

```func (*Ctd).APICreateClient(ctx *context.Context, phone string, transport string, channel_id int, nickname string, assigned_phone string) (*ClientResponse, error)```

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

```func (*Ctd).GetClientsList(ctx context.Context, offset int, limit int) (*[]ClientItem, int, error)```

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
  - A pointer to a slice of ClientItem containing the clients.
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

## Custom client fields

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

```func (*Ctd).GetCustomClientFields(ctx context.Context) (*[]CustomClientFieldItem, error)```

<details>
<summary>Function description</summary>

GetCustomClientFields retrieves a list of custom client fields from the Chat2Desk API.
It uses the APICustomClientFields method to fetch the custom client fields and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a slice of CustomClientFieldItem, which contains the custom client fields.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to a slice of CustomClientFieldItem containing the custom client fields.
</details>

</details>

## Messages

<details>
<summary>Functions list</summary>

```func (*Ctd).APISendMessage(ctx context.Context, message *MessagePayload) (*MessageResponse, error)```

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

```func (*Ctd).SendMessage(ctx context.Context, message *MessagePayload) (*MessageItem, error)```

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
  - A pointer to a slice of OperatorGroup containing the list of operator groups.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Operators

<details>
<summary>Functions list</summary>

```func (*Ctd).APIOperators(ctx context.Context, offset int, limit int) (*CtdOperatorsResponse, error)```

<details>
<summary>Function description</summary>

APIOperators retrieves a list of operators from the Chat2Desk API.
It constructs the API endpoint URL with the provided offset and limit,
sends a GET request to the API, and returns the response data as a CtdOperatorsResponse struct.
If an error occurs during the request, it logs the error and returns it.
If the request is successful, it returns a pointer to the CtdOperatorsResponse struct.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching operators.
  - limit: The maximum number of operators to return.

Returns:
  - A pointer to a CtdOperatorsResponse struct containing the list of operators and metadata.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).Operators(ctx context.Context, offset int, limit int) (*[]CtdOperator, error)```

<details>
<summary>Function description</summary>

Operators retrieves a list of operators from the Chat2Desk API.
It uses the APIOperators method to fetch the operators and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a slice of CtdOperator, which contains the operators.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination, indicating where to start fetching operators.
  - limit: The maximum number of operators to return.

Returns:
  - A pointer to a slice of CtdOperator containing the list of operators.
  - An error if the request fails or if the response is invalid.
</details>

```func (*Ctd).AllOperators(ctx context.Context) (*[]CtdOperator, error)```

<details>
<summary>Function description</summary>

AllOperators retrieves all operators from the Chat2Desk API by handling pagination.
It repeatedly calls the Operators method with increasing offsets until all operators are fetched.
It returns a pointer to a slice of CtdOperator, which contains all the operators.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to a slice of CtdOperator containing all the operators.
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

```func (*Ctd).GetTags(ctx context.Context, offset int, limit int) (*[]TagItem, int, error)```

<details>
<summary>Function description</summary>

GetTags retrieves a list of tags from the Chat2Desk API.
It uses the APIGetTags method to fetch the tags and handles errors.
If the response status is not "success", it returns nil.
It returns a pointer to a slice of TagItem, which contains the tags.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.
  - offset: The offset for pagination.
  - limit: The maximum number of tags to retrieve.

Returns:
  - A pointer to a slice of TagItem, which contains the tags
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

```func (*Ctd).GetAllTags(ctx context.Context) (*[]TagItem, error)```

<details>
<summary>Function description</summary>

GetAllTags retrieves all tags from the Chat2Desk API.
It uses the GetTags method to fetch tags in a loop until all tags are retrieved.
It returns a pointer to a slice of TagItem, which contains all the tags.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to a slice of TagItem, which contains all the tags.
  - An error if the request fails or if the response is invalid.
</details>

</details>

## Webhooks

<details>
<summary>Functions list</summary>

```func (*CreateWebhookResponse).Error() string```

<details>
<summary>Function description</summary>

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

GetWebhooks retrieves a list of webhooks from the Chat2Desk API.
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

```func (*Ctd).GetWebhooks(ctx context.Context) (*[]WebhookItem, error)```

<details>
<summary>Function description</summary>

GetWebhooks retrieves a list of webhooks from the Chat2Desk API.
It takes a context as a parameter and calls the Webhooks method.
If the response status is not "success", it logs an error and returns nil.
It returns a pointer to a slice of WebhookItem, which contains the webhooks.
If an error occurs during the request, it returns nil and the error.
If the request is successful, it returns a pointer to a slice of WebhookItem.
This method is typically used to fetch webhooks from the Chat2Desk API.

Parameters:
  - ctx: The context for the request, allowing for cancellation and timeouts.

Returns:
  - A pointer to a slice of WebhookItem containing the webhooks.
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