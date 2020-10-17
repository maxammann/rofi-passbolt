package main

// #cgo pkg-config: glib-2.0 rofi
// #include <rofi/mode.h>
// #include <rofi/mode-private.h>
// #include <gmodule.h>
// extern int rofi_init(Mode* sw);
// extern unsigned int rofi_get_num_entries(const Mode* sw);
// extern ModeMode rofi_result(Mode* sw, int menu_entry, char** input, unsigned int selected_line);
// extern void rofi_destroy(Mode* sw);
// extern int rofi_token_match(const Mode* sw, rofi_int_matcher** tokens, unsigned int index);
// extern char *rofi_get_display_value(const Mode* sw, unsigned int selected_line, int* state, GList** attr_list, int get_entry);
// extern char *rofi_get_message(const Mode *sw);
// extern char *rofi_preprocess_input(Mode* sw, const char* input);
// G_MODULE_EXPORT Mode mode;
// Mode mode =
// {
//    .abi_version        = ABI_VERSION,
//    .name               = "passbolt",
//    .cfg_name_key       = "passbolt",
//   ._init              = rofi_init,
//   ._get_num_entries   = rofi_get_num_entries,
//   ._result            = rofi_result,
//   ._destroy           = rofi_destroy,
//   ._token_match       = rofi_token_match,
//   ._get_display_value = rofi_get_display_value,
//   ._get_completion    = NULL,
//   ._preprocess_input  = NULL,
//   .private_data       = NULL,
//   .free               = NULL
// };
import "C"
