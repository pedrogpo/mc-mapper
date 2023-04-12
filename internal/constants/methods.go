package constants

import "github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"

var MethodsToMap = map[string][]string{
	"Container":    {"getSlot"},
	"Vec3":         {"addVector", "distanceTo"},
	"MouseHelper":  {"ungrabMouseCursor", "grabMouseCursor"},
	"Slot":         {"getStack", "getHasStack"},
	"ItemStack":    {"getMetadata", "getMaxDamage", "getItemUseAction"}, //  "getUnlocalizedName"
	"Item":         {"getIdFromItem"},
	"GameSettings": {"isKeyDown"},
	"KeyBinding":   {"unpressKey"},
	// "IChatComponent": {
	// 	"getFormattedText"
	// },
	"Minecraft":    {"getMinecraft", "getNetHandler"},
	"EntityPlayer": {"attackTargetEntityWithCurrentItem", "closeScreen", "getItemInUseCount"},
	"Vec3i":        {"getX", "getY", "getZ"},
	"WorldClient":  {"removeEntityFromWorld"},
	// "IWorldNameable": {
	// 	"getDisplayName"
	// },
	"AxisAlignedBB": {"expand",
		// "calculateIntercept",
		"isVecInside", "<init>"},
	"TextureManager":  {"deleteTexture", "loadTexture", "bindTexture"},
	"CapeUtils":       {"reloadCape"},
	"TileEntity":      {"getPos"},
	"CapeImageBuffer": {"<init>"},
	"Config":          {"updateFramebufferSize"},
	"World":           {"isAirBlock", "getEntityByID"},
	"Entity": {"isEntityAlive", "isInvisible", "isSprinting", "isSneaking", "setSprinting",
		// "getName",
		//  "getDisplayName",
		"getDistanceToEntity", "canBeCollidedWith", "getCollisionBorderSize", "hashCode",
		//  "getLook", "getPositionEyes"
	},
	"BlockPos":                {"<init>"},
	"Gui":                     {"drawScaledCustomSizeModalRect"},
	"IInventory":              {"getSizeInventory"},
	"RenderItem":              {"renderItemAndEffectIntoGUI"},
	"EnchantmentHelper":       {"getEnchantmentLevel"},
	"EntityLivingBase":        {"getHeldItem", "isOnSameTeam", "getHealth", "getSwingProgress", "getTeam", "getCurrentArmor"},
	"PlayerControllerMP":      {"sendUseItem", "windowClick"},
	"ResourceLocation":        {"<init>"},
	"ThreadDownloadImageData": {"setBufferedImage", "<init>"},
	"S18PacketEntityTeleport": {"getEntityId", "func_149451_c"},
	"AbstractClientPlayer":    {"getLocationSkin"},
}

func GetMethodsToMap(mappings Mappings) map[string](map[string]MethodMap) {
	methods := make(map[string](map[string]MethodMap))

	for clsName, clsMap := range mappings.Methods {
		if _, ok := MethodsToMap[clsName]; !ok {
			continue
		}

		for methodName, methodMap := range clsMap {
			find := generics.Find(MethodsToMap[clsName], func(e string) bool {
				found := false
				for _, v := range methodMap.SrgMappings {
					if v.Name == e {
						found = true
					}
				}

				for _, v := range methodMap.ObfMappings {
					if v.Name == e {
						found = true
					}
				}

				if e == methodName {
					found = true
				}
				return found
			})

			if find == nil {
				continue
			}

			// append to map
			if _, ok := methods[clsName]; !ok {
				methods[clsName] = make(map[string]MethodMap)
			}

			methods[clsName][methodName] = methodMap
		}
	}

	return methods
}

func GetMethodsToMapInClass(mappings Mappings, className string) map[string]MethodMap {
	methods := make(map[string]MethodMap)

	for clsName, clsMap := range mappings.Methods {
		if clsName != className {
			continue
		}

		for methodName, methodMap := range clsMap {
			find := generics.Find(MethodsToMap[methodMap.clsFromName], func(e string) bool {
				found := false
				for _, v := range methodMap.SrgMappings {
					if v.Name == e {
						found = true
					}
				}

				for _, v := range methodMap.ObfMappings {
					if v.Name == e {
						found = true
					}
				}

				if e == methodName {
					found = true
				}
				return found
			})

			if find == nil {
				continue
			}

			methods[methodName] = methodMap
		}
	}

	return methods

}
